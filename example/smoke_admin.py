"""goravel-admin 接口冒烟测试（对齐 owl-admin 前端契约）。

用法: BASE=http://127.0.0.1:3012 python smoke_admin.py
"""

import json
import os
import sys
import urllib.error
import urllib.parse
import urllib.request
import uuid

BASE = os.environ.get("BASE", "http://127.0.0.1:3011")
USERNAME = os.environ.get("ADMIN_USER", "smoke")
PASSWORD = os.environ.get("ADMIN_PASS", "Smoke@12345")
token = None
failed = 0

# 每次运行换一批唯一标识，脚本才能在同一套库上重复执行
RUN = uuid.uuid4().hex[:6]
PAGE_SIGN = "smoke_page_" + RUN
API_PATH = "/smoke/test_" + RUN
REL_TITLE = "Smoke Roles " + RUN
EXT_NAME = "acme/smoke-" + RUN


def call(method, path, body=None):
    url = BASE + path
    data = None
    headers = {"Accept": "application/json"}
    if body is not None:
        data = json.dumps(body).encode()
        headers["Content-Type"] = "application/json"
    if token:
        headers["Authorization"] = "Bearer " + token
    req = urllib.request.Request(url, data=data, headers=headers, method=method)
    try:
        with urllib.request.urlopen(req) as r:
            return r.status, json.loads(r.read().decode())
    except urllib.error.HTTPError as e:
        raw = e.read()
        try:
            return e.code, json.loads(raw.decode())
        except Exception:
            return e.code, {"raw": raw[:200]}


def upload(path, fields, filename, payload):
    boundary = "----smoke"
    body = b""
    for k, v in fields.items():
        body += ("--%s\r\nContent-Disposition: form-data; name=\"%s\"\r\n\r\n%s\r\n" % (boundary, k, v)).encode()
    body += ("--%s\r\nContent-Disposition: form-data; name=\"file\"; filename=\"%s\"\r\n"
             "Content-Type: application/octet-stream\r\n\r\n" % (boundary, filename)).encode()
    body += payload + ("\r\n--%s--\r\n" % boundary).encode()
    req = urllib.request.Request(BASE + path, data=body, method="POST", headers={
        "Content-Type": "multipart/form-data; boundary=" + boundary,
        "Authorization": "Bearer " + token,
    })
    with urllib.request.urlopen(req) as r:
        return json.loads(r.read().decode())


def check(name, ok, extra=""):
    global failed
    if not ok:
        failed += 1
    print(("PASS " if ok else "FAIL ") + name + ("  " + extra if extra else ""))


# 1 登录
st, r = call("POST", "/admin-api/login", {"username": USERNAME, "password": PASSWORD})
token = (r.get("data") or {}).get("token")
check("login", st == 200 and r.get("status") == 0 and bool(token), "code=%s" % r.get("status"))

# 1.1 匿名必须被拦住（用户管理曾因白名单配置漏在外网）
saved, token = token, None
st, r = call("GET", "/admin-api/system/admin_users?_action=getData")
token = saved
check("anonymous blocked", r.get("code") == 401, "code=%s" % r.get("code"))

# 1.2 用户列表分页字段：amis 靠 total 才能渲染分页条
st, r = call("GET", "/admin-api/system/admin_users?_action=getData&page=1&perPage=1")
d = r.get("data") or {}
check("user list paginated", st == 200 and "total" in d and len(d.get("items") or []) == 1,
      "total=%s items=%s" % (d.get("total"), len(d.get("items") or [])))

# 1.3 关键字过滤要真的生效（Where 未赋值会让条件被静默丢弃）
st, r = call("GET", "/admin-api/system/admin_users?_action=getData&keyword=zzz_not_exist")
check("user list keyword filter", st == 200 and not ((r.get("data") or {}).get("items")), "")

# 2 站点设置
st, r = call("GET", "/admin-api/_settings")
d = r.get("data") or {}
check("_settings", st == 200 and "enabled_extensions" in d and "assets" in d,
      "enabled_extensions=%s" % d.get("enabled_extensions"))

# 3 验证码
st, r = call("GET", "/admin-api/captcha")
img = (r.get("data") or {}).get("captcha_img", "")
check("captcha", st == 200 and img.startswith("data:image/"), "bytes=%d" % len(img))

# 4 iconify
st, r = call("GET", "/admin-api/_iconify_search?keyword=home")
items = (r.get("data") or {}).get("items") or []
check("_iconify_search", st == 200 and len(items) > 0, "count=%d" % len(items))

# 5 自定义页面（前端契约：page.schema 嵌套）
st, r = call("DELETE", "/admin-api/dev_tools/pages/1")
st, r = call("POST", "/admin-api/dev_tools/pages", {
    "title": "Smoke Page " + RUN, "sign": PAGE_SIGN,
    "page": {"schema": {"type": "page", "body": [{"type": "tpl", "tpl": "hello"}]}},
})
created = r.get("status") == 0
st, r2 = call("POST", "/admin-api/dev_tools/pages", {
    "title": "Dup", "sign": PAGE_SIGN, "page": {"schema": {"type": "page"}},
})
dup_rejected = r2.get("status") != 0
st, r3 = call("GET", "/admin-api/page_schema?sign=" + PAGE_SIGN)
schema = r3.get("data")
check("pages store/unique/page_schema",
      created and dup_rejected and isinstance(schema, dict) and schema.get("type") == "page",
      "create=%s dup=%s schema_type=%s" % (r.get("status"), r2.get("status"), type(schema).__name__))

# 6 页面更新：非法 id 不得插新行
st, r = call("PUT", "/admin-api/dev_tools/pages/999999", {
    "id": 999999, "title": "Ghost", "sign": "ghost_page", "page": {"schema": {"type": "page"}},
})
st, rl = call("GET", "/admin-api/dev_tools/pages?_action=getData&page=1&perPage=50")
titles = [i.get("title") for i in (rl.get("data") or {}).get("items", [])]
check("pages update ghost id rejected", r.get("status") != 0 and "Ghost" not in titles,
      "status=%s rows=%s" % (r.get("status"), titles))

# 7 动态接口
api = {"title": "Smoke Api", "path": API_PATH, "method": "get", "template": "list",
       "enabled": True, "args": {"model": "admin_users"}}
st, a1 = call("POST", "/admin-api/dev_tools/apis", api)
st, a2 = call("POST", "/admin-api/dev_tools/apis", api)
st, tpl = call("GET", "/admin-api/dev_tools/api/templates")
check("apis store/dup/templates",
      a1.get("status") == 0 and a2.get("status") != 0 and len(tpl.get("data") or []) >= 10,
      "first=%s dup=%s templates=%d" % (a1.get("status"), a2.get("status"), len(tpl.get("data") or [])))

# 8 关联
st, r = call("GET", "/admin-api/dev_tools/relation/types")
types = r.get("data") or []
st, rr = call("POST", "/admin-api/dev_tools/relationships", {
    "model": "admin_users", "title": REL_TITLE, "type": "belongsToMany",
    "args": {"related": "admin_roles", "table": "role_user", "foreign_key": "user_id", "related_key": "role_id"},
})
st, rlist = call("GET", "/admin-api/dev_tools/relationships?_action=getData&page=1&perPage=20")
rows = (rlist.get("data") or {}).get("items", [])
code = next((i.get("preview_code", "") for i in rows if i.get("title") == REL_TITLE), "")
check("relationships types/preview",
      len(types) >= 9 and rr.get("status") == 0
      and "func (m *AdminUsers) SmokeRoles" in code and "role_user" in code,
      "types=%d create=%s" % (len(types), rr.get("status")))

# 9 扩展：名单化 + 启用；非法 id 必须拒绝
st, bad = call("POST", "/admin-api/dev_tools/extensions", {"name": "badname"})
st, ok_ext = call("POST", "/admin-api/dev_tools/extensions", {"name": EXT_NAME, "title": "Demo " + RUN})
st, elist = call("GET", "/admin-api/dev_tools/extensions?_action=getData&page=1&perPage=50")
ext = next((i for i in (elist.get("data") or {}).get("items", []) if i.get("name") == EXT_NAME), {})
st, en = call("POST", "/admin-api/dev_tools/extensions/enable", {"id": ext.get("id")})
st, en_bad = call("POST", "/admin-api/dev_tools/extensions/enable", {"id": 0})
st, s2 = call("GET", "/admin-api/_settings")
enabled = (s2.get("data") or {}).get("enabled_extensions") or []
check("extensions store/enable/guard",
      bad.get("status") != 0 and ok_ext.get("status") == 0 and en.get("status") == 0
      and en_bad.get("status") != 0 and EXT_NAME in enabled and "" not in enabled,
      "bad=%s ok=%s enable=%s guard=%s enabled=%s" % (
          bad.get("status"), ok_ext.get("status"), en.get("status"), en_bad.get("status"), enabled))

# 10 分块上传
st, r = call("POST", "/admin-api/upload_chunk_start",
             {"filename": "chunk.bin", "size": 1536, "chunkSize": 1024, "disk": "admin"})
uid = (r.get("data") or {}).get("uploadId")
p1 = upload("/admin-api/upload_chunk", {"uploadId": uid, "partNumber": "1", "size": "1024"}, "blob", b"A" * 1024)
p2 = upload("/admin-api/upload_chunk", {"uploadId": uid, "partNumber": "2", "size": "512"}, "blob", b"B" * 512)
st, rfin = call("POST", "/admin-api/upload_chunk_finish", {
    "uploadId": uid, "filename": "chunk.bin", "disk": "admin",
    "partList": [{"partNumber": 1, "eTag": p1["data"]["eTag"]}, {"partNumber": 2, "eTag": p2["data"]["eTag"]}]})
check("chunk upload", bool(uid) and p1.get("status") == 0 and p2.get("status") == 0
      and rfin.get("status") == 0, "path=%s" % (rfin.get("data") or {}).get("path"))

# 11 导出
st, r = call("GET", "/admin-api/system/admin_users?_action=export")
d = r.get("data") or {}
ok_sign = all(k in d for k in ("path", "expires", "signature"))
if ok_sign:
    q = urllib.parse.urlencode({"path": d["path"], "expires": d["expires"], "signature": d["signature"]})
    with urllib.request.urlopen(urllib.request.Request(BASE + "/admin-api/_download_export?" + q)) as resp:
        content = resp.read()
    bad_q = urllib.parse.urlencode(dict(d, signature="0" * 64))
    try:
        with urllib.request.urlopen(BASE + "/admin-api/_download_export?" + bad_q) as resp:
            tampered = resp.status
    except urllib.error.HTTPError as e:
        tampered = e.code
    check("export sign/download/tamper",
          content[:2] == b"PK" and tampered in (401, 403),
          "bytes=%d tampered=%s" % (len(content), tampered))
else:
    check("export sign/download/tamper", False, "missing keys: %s" % list(d.keys()))

# 12 个人设置
st, r = call("PUT", "/admin-api/user_setting", {"name": "Smoke Renamed", "email": "smoke@example.com"})
st, r2 = call("GET", "/admin-api/current-user")
check("user_setting", r.get("status") == 0 and (r2.get("data") or {}).get("name") == "Smoke Renamed",
      "name=%s" % (r2.get("data") or {}).get("name"))

# 13 菜单
st, r = call("GET", "/admin-api/menus")
check("menus", st == 200 and len(r.get("data") or []) > 0, "count=%d" % len(r.get("data") or []))

# 14 列表首页正常渲染
for path in ("/admin-api/system/admin_users", "/admin-api/system/admin_roles",
             "/admin-api/system/admin_menus", "/admin-api/system/admin_permissions"):
    st, r = call("GET", path)
    ok = st == 200 and isinstance(r.get("data"), dict) and "body" in (r.get("data") or {})
    check("page schema %s" % path, ok, "keys=%s" % list((r.get("data") or {}).keys())[:6])

print("\n%d 项失败" % failed)
sys.exit(1 if failed else 0)

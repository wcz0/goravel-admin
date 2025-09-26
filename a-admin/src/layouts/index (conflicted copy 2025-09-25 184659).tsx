import React from 'react'
import useRoute from '@/routes'
import DefaultLayout from '@/layouts/DefaultLayout'
import LayoutContent from '@/layouts/components/LayoutContent'
import SmLayout from '@/layouts/SmLayout'
import useSmallScreen from '@/hooks/useSmallScreen'
import TopLayout from '@/layouts/TopLayout'
import TopMixLayout from '@/layouts/TopMixLayout'
import {DoubleLayout} from '@/layouts/DoubleLayout'
import SettingPanel from '@/layouts/components/SettingPanel'
import useSetting from '@/hooks/useSetting'
import {goToLoginPage, Token} from '@/utils/common'

// 应用整体布局
const Layout = () => {
    const {getCurrentRoute} = useRoute()
    const isSmallScreen = useSmallScreen()
    const {getSetting} = useSetting()

    const currentRoute = getCurrentRoute()

    // 根据配置渲染不同布局
    const RenderLayout = (mode: string) => {
        switch (mode) {
            case 'default':
                return <DefaultLayout/>
            case 'top':
                return <TopLayout/>
            case 'top-mix':
                return <TopMixLayout/>
            case 'double':
                return <DoubleLayout/>
            default:
                return <></>
        }
    }

    // 判断是否登录
    if (!Token().value) {
        goToLoginPage()
    }

    // 判断是否全屏
    if (currentRoute?.is_full == 1) {
        return (
            <div className="h-screen">
                <LayoutContent/>
            </div>
        )
    }

    return (
        <>
            {isSmallScreen ? <SmLayout/> : RenderLayout(getSetting('system_theme_setting.layoutMode', 'default'))}
            <SettingPanel/>
        </>
    )
}

export default Layout

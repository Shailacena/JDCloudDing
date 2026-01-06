import { Route, Routes } from 'react-router-dom'
import './App.css'
import AppProvider, { RequireAuth, RequireAuthMerchant, RequireAuthPartner } from './AppProvider'
import MainLayout from './admin/MainLayout'
import Login from './admin/Login'
import LoginPartner from './partner/Login'
import LoginMerchant from './merchant/Login'
import MainLayoutPartner from './partner/MainLayout'
import MainLayoutMerchant from './merchant/MainLayout'
import { getRouteConfig as getRouteConfigPartner } from './partner/RouteConfigs'
import { getRouteConfig as getRouteConfigMerchant } from './merchant/RouteConfigs'
import { IRoute, routes } from './admin/routes'
import zhCN from 'antd/lib/locale/zh_CN';
import { AxiosProvider } from './api/AxiosProvider'
import { ConfigProvider } from 'antd'
import PaymentPage from './open/PaymentPage'

// import MainLayout from './MainLayout'

function App() {
  return (
    // <>
    //   <MainLayout />
    // </>
    <ConfigProvider locale={zhCN}>
      <AppProvider>
        <AxiosProvider>

          <Routes>
            <Route path="/admin/login" element={<Login />} />
            <Route path="/partner/login" element={<LoginPartner />} />
            <Route path="/merchant/login" element={<LoginMerchant />} />
            <Route path="/open/PaymentPage" element={<PaymentPage />} />
            <Route
              path="/admin"
              element={
                <RequireAuth>
                  <MainLayout />
                </RequireAuth>
              }
            >
              {routes.map((r: IRoute) => {
                let route = r?.children ? r.children.map((childRoute: IRoute) => (
                  <Route key={r.path + childRoute.path} path={r.path + childRoute.path} element={<childRoute.component />} />
                )) : <Route key={r.path} path={r.path} element={<r.component />} />
                return route
              })}
            </Route>

            <Route
              path="/partner"
              element={
                <RequireAuthPartner>
                  <MainLayoutPartner />
                </RequireAuthPartner>
              }
            >
              {getRouteConfigPartner().map((r: IRoute) => {
                let route = r?.children ? r.children.map((childRoute: IRoute) => (
                  <Route key={r.path + childRoute.path} path={r.path + childRoute.path} element={<childRoute.component />} />
                )) : <Route key={r.path} path={r.path} element={<r.component />} />
                return route
              })}
            </Route>

            <Route
              path="/merchant"
              element={
                <RequireAuthMerchant>
                  <MainLayoutMerchant />
                </RequireAuthMerchant>
              }
            >
              {getRouteConfigMerchant().map((r: IRoute) => {
                let route = r?.children ? r.children.map((childRoute: IRoute) => (
                  <Route key={r.path + childRoute.path} path={r.path + childRoute.path} element={<childRoute.component />} />
                )) : <Route key={r.path} path={r.path} element={<r.component />} />
                return route
              })}
            </Route>
          </Routes>
        </AxiosProvider>
      </AppProvider>
    </ConfigProvider>
  )
}

export default App

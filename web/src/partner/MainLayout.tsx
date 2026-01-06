import { Layout, Menu, Dropdown, Col, Row, Flex, message } from 'antd';
import type { MenuProps } from 'antd';
import { DownOutlined } from '@ant-design/icons';
import { getRouteConfig } from './RouteConfigs';
import { Outlet, useLocation, useNavigate } from 'react-router-dom';
import { useAppContext } from '../AppProvider';
import { useState } from 'react';
import SetPasswordModal from './modal/SetPasswordModal';
import { useApis } from '../api/api';

type MenuItem = Required<MenuProps>['items'][number];

const { Header, Content, Sider, Footer } = Layout;

const menus: MenuItem[] = getRouteConfig().map(
  (menu) => {
    return {
      key: menu.path,
      icon: menu.icon ? <menu.icon /> : null,
      label: menu.name,

      children: menu.children?.map((sub) => {
        return {
          key: menu.path + sub.path,
          label: sub.name,
        };
      }),
    };
  },
);

enum DropdownMenuKey {
  SetPassword = '0',
  Logout = '1'
}

const items: MenuProps['items'] = [
  {
    label: "修改密码",
    key: DropdownMenuKey.SetPassword
  },
  {
    label: "登出",
    key: DropdownMenuKey.Logout
  }
]

function MainLayout() {
  const navigate = useNavigate()
  const loc = useLocation()
  const ctx = useAppContext()
  const { partner1Logout } = useApis()
  const [isOpenPasswordModal, setOpenPasswordModal] = useState(false)

  const onClickMenu: MenuProps['onClick'] = (e) => {
    if (loc.pathname == e.key) {
      return
    }
    navigate(e.key)
  }

  const logout = async () => {
    try {
      await partner1Logout({})

      message.success("退出登录成功", 0.7, () => {
        ctx.auth.partnerSignout(() => { })
      })
    } catch (error) {
      console.log(error);
    }
  }

  const dropdownMenuClick: MenuProps['onClick'] = ({ key }) => {
    switch (key) {
      case DropdownMenuKey.SetPassword:
        setOpenPasswordModal(true)
        break;
      case DropdownMenuKey.Logout:
        logout()
        break;
    }
  }
  return (
    <>
      {
        isOpenPasswordModal &&
        <SetPasswordModal isModalOpen={isOpenPasswordModal} onOk={() => { setOpenPasswordModal(false) }} onCancel={() => { setOpenPasswordModal(false) }} />
      }
      <Layout>
        <Header style={{ color: "#fff", height: 48 }}>
          <Row style={{ height: "100%" }}>
            <Col style={{ height: "100%", lineHeight: "48px" }} span={12}>
              <span style={{ fontSize: "25px" }}>管理后台</span>
            </Col>
            <Col style={{ height: "100%", lineHeight: "48px" }} span={12}>
              <Flex justify="center" align="end" vertical>
                <Dropdown menu={{
                  items,
                  onClick: dropdownMenuClick
                }} trigger={['click']}>
                  <a style={{ color: "#fff", fontSize: "16px" }} onClick={(e) => e.preventDefault()}>
                    {ctx?.cookie?.nickname}
                    <DownOutlined />
                  </a>
                </Dropdown>
              </Flex>
            </Col>
          </Row>
        </Header>
        <Layout>
          <Sider width={200}>
            <Menu
              mode="inline"
              style={{ height: '100%' }}
              items={menus}
              onClick={onClickMenu}
            />
          </Sider>
          <Content style={{ padding: '20px', height: "calc(100vh - (40px + 48px))" }}>
            <Outlet />
          </Content>
        </Layout>
        <Footer style={{ padding: "12px 50px", textAlign: 'center' }}>
          Copyright ©{new Date().getFullYear()} Content-Manage-System
        </Footer>
      </Layout >
    </>
  )
}

export default MainLayout

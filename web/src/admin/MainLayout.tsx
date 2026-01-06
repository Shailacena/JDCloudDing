import { Layout, Menu, Dropdown, message, Row, Col, Flex } from 'antd';
import type { MenuProps } from 'antd';
import { DownOutlined } from '@ant-design/icons';
import { Outlet, useLocation, useNavigate } from 'react-router-dom'
import { routes } from './routes';
import { useAppContext } from '../AppProvider';
import { useEffect, useState } from 'react';
import CurrentLocation from '../components/CurrentLocation';
import { useApis } from '../api/api';
import SetPasswordModal from './modal/SetPasswordModal';
import { RoleType } from './role';

const { Header, Content, Sider, Footer } = Layout;

type MenuItem = Required<MenuProps>['items'][number];

enum DropdownMenuKey {
  SetPassword = '0',
  Logout = '1'
}

interface LevelKeysProps {
  key?: string;
  children?: LevelKeysProps[];
}

const getLevelKeys = (items: LevelKeysProps[]) => {
  const key: Record<string, number> = {};
  const func = (items2: LevelKeysProps[], level = 1) => {
    items2.forEach((item) => {
      if (item.key) {
        key[item.key] = level;
      }
      if (item.children) {
        func(item.children, level + 1);
      }
    });
  };
  func(items);
  return key;
};

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
  const [menuItems, setMenuItems] = useState<MenuItem[]>([])
  const navigate = useNavigate()
  const loc = useLocation()
  const ctx = useAppContext()
  const { adminLogout } = useApis()
  const [isOpenPasswordModal, setOpenPasswordModal] = useState(false)
  const [defaultSelectedKeys, setDefaultSelectedKeys] = useState(loc.pathname)

  const nickname = ctx?.cookie?.nickname

  useEffect(() => {
    let all = routes.filter((route) => {
      if ((ctx.cookie.role === RoleType.Agency || ctx.cookie.role === RoleType.ClonedAdmin) && (route.inPermission && ctx.cookie.role > route.inPermission)) {
        return false
      }

      let isShow = route.permission ? route.permission === ctx.cookie.role : true
      if (!isShow) {
        return false
      }

      route.children = route.children?.filter((subRoute) => {
        return subRoute.permission ? subRoute.permission === ctx.cookie.role : true
      })
      return true
    })

    let menu = all.map(
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

    setMenuItems(menu)
  }, [ctx.cookie])

  useEffect(() => {
    setDefaultSelectedKeys(loc.pathname)
  }, [loc.pathname])

  const onClickMenu: MenuProps['onClick'] = (e) => {
    if (loc.pathname == e.key) {
      return
    }
    navigate(e.key)
  }

  const levelKeys = getLevelKeys(menuItems as LevelKeysProps[]);
  const [stateOpenKeys, setStateOpenKeys] = useState([routes[0].path, routes[0].path]);

  const onOpenChange: MenuProps['onOpenChange'] = (openKeys) => {
    const currentOpenKey = openKeys.find((key) => stateOpenKeys.indexOf(key) === -1);
    if (currentOpenKey !== undefined) {
      const repeatIndex = openKeys
        .filter((key) => key !== currentOpenKey)
        .findIndex((key) => levelKeys[key] === levelKeys[currentOpenKey]);
      setStateOpenKeys(
        [openKeys
          .filter((_, index) => index !== repeatIndex)
          .filter((key) => levelKeys[key] <= levelKeys[currentOpenKey]).pop() || openKeys[0]],
      );
    } else {
      setStateOpenKeys(openKeys);
    }
  }

  const logout = async () => {
    await adminLogout({})

    message.success("退出登录成功", 0.7, () => {
      ctx.auth.adminSignout(() => { })
    })
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
                    {nickname}
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
              defaultSelectedKeys={[defaultSelectedKeys]}
              items={menuItems}
              openKeys={stateOpenKeys}
              onOpenChange={onOpenChange}
              onClick={onClickMenu}
            />
          </Sider>
          <Content style={{ padding: '20px', height: "calc(100vh - (40px + 48px))", overflow: "overlay" }}>
            <div style={{ marginBottom: '10px' }}>
              <CurrentLocation routeconfigs={routes} />
            </div>
            <Outlet />
          </Content>
        </Layout>
        <Footer style={{ padding: "12px 50px", textAlign: 'center' }}>
          Copyright ©{new Date().getFullYear()} 管理后台
        </Footer>
      </Layout>

      {
        isOpenPasswordModal &&
        <SetPasswordModal isModalOpen={isOpenPasswordModal} onOk={() => { setOpenPasswordModal(false) }} onCancel={() => { setOpenPasswordModal(false) }} />
      }
    </>
  )
}

export default MainLayout

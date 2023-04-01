import React, { useEffect, useState } from 'react';
import { Link, Outlet, useLocation } from 'react-router-dom'
import { observer } from 'mobx-react-lite';
import { Avatar, Button, Col, Dropdown, Layout, Menu, Row } from 'antd';
import { ItemType } from 'antd/es/menu/hooks/useItems';
import { AiOutlineCaretDown, AiOutlineUser } from 'react-icons/ai';

import { fixMenuSelected } from 'transforms/route';
import { useGlobalContext } from 'global/Store';

// fix menu selected
const patterns: { [props: string]: string } = {
  'spaces/:key': '/spaces',
}

const Admin = observer(() => {

  const location = useLocation();
  const { account } = useGlobalContext();

  const [accountDropdownMenuItems, setAccountDropdownMenu] = useState<ItemType[]>([])

  const menuItems: ItemType[] = [
    {
      key: '/spaces',
      label: <Link to="/spaces">Spaces</Link>
    },
    {
      key: '/accounts',
      label: <Link to="/accounts">Accounts</Link>
    },
  ];

  const [menuSelectedKeys, setMenuSelectedKeys] = useState<string[]>([]);

  // useEffect(() => {
  //   const path = location.pathname.split('/');
  //   path.shift();
  //   const match = path.map((item, index, array) => index === 0 ? `/${item}` : `/${array.slice(0, index).join('/')}/${item}`);
  //   setMenuSelectedKeys([location.pathname, ...match]);
  // }, [location.pathname]);

  useEffect(() => {
    const items: ItemType[] = [
      { key: 'current_user', label: <Link to='/spaces'>Signed in as <strong>{account.name || account.login}</strong></Link> },
      { type: 'divider' },
      { key: 'sign_out', label: <a href="/logout">Sign out</a> },
    ]
    setAccountDropdownMenu(items)
  }, [account])

  // change menu selected when route changed
  useEffect(() => {
    const match = fixMenuSelected(patterns, location.pathname)
    setMenuSelectedKeys([match]);
  }, [location.pathname]);

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Layout.Header
        className="app-layout-header"
        style={{
          position: 'fixed',
          top: 0,
          zIndex: 9,
          width: '100%'
        }}
      >
        <Row>
          <Col flex="100px">
            <Link to="/spaces" className="logo" style={{ display: 'inline-block' }}>
              <svg width="80" height="32" viewBox="0 0 80 22" fill="none">
                <text fill="#FFF" style={{ fontWeight: 'bold' }} x="5" y="18" fontFamily="Verdana" fontSize="20">Space</text>
              </svg>
            </Link>
          </Col>
          <Col flex="auto">
            <Menu
              className="navigation"
              theme="dark"
              mode="horizontal"
              inlineIndent={8}
              defaultSelectedKeys={[location.pathname]}
              selectedKeys={menuSelectedKeys}
              items={menuItems}
            />
          </Col>
          <Col>
            <Dropdown
              trigger={['click']}
              placement="bottomRight"
              arrow
              menu={{ items: accountDropdownMenuItems }}
            >
              <Button type="link" style={{ paddingRight: 0 }}>
                <Avatar
                  size={20}
                  icon={<AiOutlineUser />}
                  src={account.avatar}
                  style={{ backgroundColor: '#87d068' }}
                />
                <AiOutlineCaretDown style={{ margin: 0 }} />
              </Button>
            </Dropdown>
          </Col>
        </Row>
      </Layout.Header>

      <Layout style={{ paddingTop: 64 }}>
        <Outlet />
      </Layout>
    </Layout>
  );
})

export default Admin;

import React, { useMemo } from 'react';
import { observer } from 'mobx-react-lite';
import { BrowserRouter, Link, Route, Routes } from 'react-router-dom'
import { QueryParamProvider } from 'use-query-params';
import { ReactRouter6Adapter } from 'use-query-params/adapters/react-router-6';
import { useQuery } from '@tanstack/react-query';
import queryString from 'query-string';
import { ConfigProvider, Typography } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { IconContext } from 'react-icons';

import { WaitingComponent } from 'components/WaitingComponent';
import { GlobalContext, GlobalStore } from 'global/Store';
import ApplicationLayout from 'pages/ApplicationLayout';

const Signup = WaitingComponent(React.lazy(() => import(/* webpackChunkName: "signup" */ 'pages/Signup')));
const Signin = WaitingComponent(React.lazy(() => import(/* webpackChunkName: "signin" */ 'pages/Signin')));

const Spaces = WaitingComponent(React.lazy(() => import(/* webpackChunkName: "spaces" */ 'pages/Spaces')));
const NewSpace = WaitingComponent(React.lazy(() => import(/* webpackChunkName: "spaces" */ 'pages/Spaces/New')));
const Space = WaitingComponent(React.lazy(() => import(/* webpackChunkName: "spaces" */ 'pages/Spaces/Detail')));
const SpaceDashboard = WaitingComponent(React.lazy(() => import(/* webpackChunkName: "spaces" */ 'pages/Spaces/Detail/dashboard')));
const EditSpace = WaitingComponent(React.lazy(() => import(/* webpackChunkName: "spaces" */ 'pages/Spaces/Edit')));
const Page = WaitingComponent(React.lazy(() => import(/* webpackChunkName: "spaces" */ 'pages/Spaces/Page')));
const NewPage = WaitingComponent(React.lazy(() => import(/* webpackChunkName: "spaces" */ 'pages/Pages/New')));
const EditPage = WaitingComponent(React.lazy(() => import(/* webpackChunkName: "spaces" */ 'pages/Pages/Edit')));

const Accounts = WaitingComponent(React.lazy(() => import(/* webpackChunkName: "accounts" */ 'pages/Accounts')));

const App = observer(() => {

  const store = useMemo<GlobalStore>(() => new GlobalStore(), []);

  useQuery(['overview'], () => store.loadOverview())

  return (
    <GlobalContext.Provider value={store}>
      <ConfigProvider locale={zhCN} prefixCls="app">
        <IconContext.Provider value={{ className: 'app-icon' }}>
          <BrowserRouter>
            <QueryParamProvider
              adapter={ReactRouter6Adapter}
              options={{
                searchStringToObject: queryString.parse,
                objectToSearchString: queryString.stringify,
              }}
            >
              <Routes>
                <Route path="/signup" element={<Signup />} />
                <Route path="/signin" element={<Signin />} />

                <Route path="" element={<ApplicationLayout />}>
                  {/* <Route index element={<Dashboard />} /> */}

                  <Route index element={<Spaces />} />
                  <Route path="spaces" element={<Spaces />} />
                  <Route path="spaces/new" element={<NewSpace />} />

                  <Route path="spaces/:key" element={<Space />}>
                    <Route index element={<SpaceDashboard />} />
                    <Route path="setting/profile" element={<EditSpace />} />
                    <Route path="pages/:page_id" element={<Page />} />
                    <Route path="pages/new" element={<NewPage />} />
                    <Route path="pages/:page_id/edit" element={<EditPage />} />
                  </Route>

                  <Route path="accounts" element={<Accounts />} />
                </Route>

                <Route path='/forbidden' element={<Forbidden />} />
                <Route path='/500' element={<Oops />} />
                <Route path="*" element={<NoMatch />} />

              </Routes>
            </QueryParamProvider>
          </BrowserRouter>
        </IconContext.Provider>
      </ConfigProvider>
    </GlobalContext.Provider>
  );
})

function NoMatch() {
  return (
    <div style={{ padding: 50 }}>
      <Typography.Title level={1} style={{ fontSize: 72, margin: 0 }}>404</Typography.Title>
      <Typography.Title level={2} style={{ margin: 0 }}>Page Not Found.</Typography.Title>
      <p style={{ marginTop: 15 }}>
        <Link to="/">Go to the home page</Link>
      </p>
    </div>
  );
}

function Forbidden() {
  return (
    <div style={{ padding: 50 }}>
      <Typography.Title level={1} style={{ fontSize: 72, margin: 0 }}>403</Typography.Title>
      <Typography.Title level={2} style={{ margin: 0 }}>Forbidden.</Typography.Title>
      <Typography.Paragraph style={{ marginTop: 15 }}>Contact the administrator for permission</Typography.Paragraph>
      <Typography.Paragraph style={{ marginTop: 15 }}>
        or, return to the <a href="/">home page</a>
      </Typography.Paragraph>
    </div>
  );
}

function Oops() {
  return (
    <div style={{ padding: 50 }}>
      <Typography.Title level={1} style={{ fontSize: 72, margin: 0 }}>500</Typography.Title>
      <Typography.Title level={2} style={{ margin: 0 }}>Oops, Internal Server Error.</Typography.Title>
      <Typography.Paragraph style={{ marginTop: 15 }}>Contact the administrator, or return to the <a href="/">home page</a>
      </Typography.Paragraph>
    </div>
  );
}

export default App;

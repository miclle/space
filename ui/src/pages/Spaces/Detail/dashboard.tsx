import { observer } from "mobx-react-lite";
import { PageHeader } from '@ant-design/pro-components';

import { useSpaceContext } from "./store";

const Dashboard = observer(() => {
  const store = useSpaceContext();
  const { space } = store;

  return (
    <>
      <PageHeader
        ghost={false}
        title={space.name}
      />

      <h1>Space Dashboard</h1>
      <p>TODO</p>
    </>
  );
})

export default Dashboard
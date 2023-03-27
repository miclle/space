import React, { Suspense } from 'react';
import { Skeleton } from 'antd';

/**
 * Component Waiting wrapper
 * @param Component React lazy component
 */
export function WaitingComponent(Component: React.LazyExoticComponent<(config: any) => JSX.Element>) {
  return (config: any): JSX.Element => (
    <Suspense fallback={<Skeleton active />}>
      <Component {...config} />
    </Suspense>
  );
}

export default WaitingComponent;

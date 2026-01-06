import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { Breadcrumb } from 'antd';
import { IRoute } from '../admin/routes';
import { BreadcrumbItemType, BreadcrumbSeparatorType } from 'antd/es/breadcrumb/Breadcrumb';

// 定义props的接口
interface CurrentLocationProps {
  routeconfigs: IRoute[];
}

const CurrentLocation: React.FC<CurrentLocationProps> = ({ routeconfigs }) => {

  const location = useLocation();
  // 使用路由路径找到模块名称
  let routes: IRoute[] = []
  getRoutesFormPath(routes, routeconfigs, location.pathname);

  function getRoutesFormPath(routes: IRoute[], routeconfigs: IRoute[], path: string) {
    let leftPath = path;
    for (let route of routeconfigs) {
      let index = leftPath.indexOf(route.path);
      if (index != -1) {
        routes.push(route)
        if (route.children && route.children.length > 0) {
          leftPath = leftPath.replace(route.path, '');
          getRoutesFormPath(routes, route.children, leftPath);
        }
      }
    }
  }

  function itemRender(currentRoute: any, _: any, items: any, paths: any) {
    let isLast = currentRoute?.path === items[items.length - 1]?.path;
    let isFirst = currentRoute?.path === items[0]?.path;
    return isLast ? (<span>{currentRoute.title}</span>) :
      (isFirst ? (<span>{currentRoute.title}</span>) : (<Link to={`/${paths.join("/")}`}>{currentRoute.title}</Link>));
  }

  const breadcrumbItems = routes.map(route => ({
    title: route.name,
    path: route.path
  } as Partial<BreadcrumbItemType & BreadcrumbSeparatorType>))

  return <Breadcrumb itemRender={itemRender} items={breadcrumbItems} />;
};

export default CurrentLocation;

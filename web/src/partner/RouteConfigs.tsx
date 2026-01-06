import { HomeOutlined, TableOutlined, ShopOutlined } from '@ant-design/icons';
import Dashboard from './Dashboard';
import { getRandomPath } from '../utils/Tool';
import Goods from './Goods';
import CashFlowDaily from './CashFlowDaily';
import BalanceBill from './BalanceBill';
import TradingRecord from './TradingRecord';

export interface IRoute {
  name: string
  path: string
  component: any,
  icon?: any,
  children?: Array<IRoute>
}

enum MODEL_PATH {
  HOME, ORDERS, CASHFLOW
}

const routeconfigs: IRoute[] = []

export function getRouteConfig(): Array<IRoute> {
  if (routeconfigs.length != 0) {
    return routeconfigs;
  } else {
    let routs = [{
      // path: '/partner/home',
      path: '/partner/' + getRandomPath(MODEL_PATH.HOME),
      name: '首页',
      component: Dashboard,
      icon: HomeOutlined,
    },
    {
      path: '/partner/' + getRandomPath(MODEL_PATH.ORDERS),
      name: '商品管理',
      icon: TableOutlined,
      component: null,
      children: [
        {
          path: '/' + getRandomPath(0),
          name: '商品列表',
          component: Goods,
        },
      ],
    },
    {
      path: '/partner/' + getRandomPath(MODEL_PATH.CASHFLOW),
      name: '交易管理',
      icon: ShopOutlined,
      component: null,
      children: [
        {
          path: '/' + getRandomPath(0),
          name: '交易记录',
          component: TradingRecord,
        },
        {
          path: '/' + getRandomPath(1),
          name: '余额明细',
          component: BalanceBill,
        },
        {
          path: '/' + getRandomPath(2),
          name: '每日流水',
          component: CashFlowDaily,
        },
      ],
    },
    ];
    for (let i in routs) {
      routeconfigs.push(routs[i]);
    }
    return routeconfigs;
  }
}


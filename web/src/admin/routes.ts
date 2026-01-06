import { HomeOutlined, MergeCellsOutlined, UserOutlined, ShopOutlined, PayCircleOutlined, TableOutlined } from '@ant-design/icons';
import Dashboard from './Dashboard';
import ServerStatus from './ServerStatus';
import Admin from './Admin';
import Partner from './Partner';
import Goods from './Goods';
import Merchant from './Merchant';
import TradingRecord from './TradingRecord';
import { RoleType } from './role';
import MerchantBalanceBill from './MerchantBalanceBill';
import PartnerBalanceBill from './PartnerBalanceBill';
import OperationLog from './OperationLog';
import PartnerDailyBill from './PartnerDailyBill';
import MerchantDailyBill from './MerchantDailyBill';
import DailyBill from './DailyBill';

export interface IRoute {
  name: string
  path: string
  component: any,
  icon?: any,
  children?: Array<IRoute>
  permission?: RoleType
  inPermission?: RoleType
}

export const routes: Array<IRoute> = [
  {
    path: '/admin/home',
    name: '首页',
    component: Dashboard,
    icon: HomeOutlined,
  },
  {
    path: '/admin/serverStatus',
    name: '服务器状态',
    component: ServerStatus,
    icon: MergeCellsOutlined,
  },
  {
    path: '/admin/manager',
    name: '管理员',
    icon: UserOutlined,
    component: null,
    inPermission: RoleType.Admin,
    children: [
      {
        path: '/list',
        name: '管理员列表',
        component: Admin,
      },
      {
        path: '/operationLog',
        name: '操作日志',
        component: OperationLog,
      },
    ],
  },
  {
    path: '/admin/partner',
    name: '合作商管理',
    icon: ShopOutlined,
    component: null,
    children: [
      {
        path: '/list',
        name: '合作商列表',
        component: Partner,
      },
      {
        path: '/goods',
        name: '商品管理',
        component: Goods,
      },
      {
        path: '/balanceBill',
        name: '余额明细',
        component: PartnerBalanceBill,
      },
      {
        path: '/dailyBill',
        name: '每日流水',
        component: PartnerDailyBill,
      },
    ],
  },
  {
    path: '/admin/merchant',
    name: '商户管理',
    icon: MergeCellsOutlined,
    component: null,
    children: [
      {
        path: '/list',
        name: '商户列表',
        component: Merchant,
      },
      {
        path: '/balanceBill',
        name: '余额明细',
        component: MerchantBalanceBill,
      },
      {
        path: '/dailyBill',
        name: '每日流水',
        component: MerchantDailyBill,
      },
    ],
  },
  {
    path: '/admin/trade',
    name: '交易统计',
    icon: PayCircleOutlined,
    component: null,
    children: [
      {
        path: '/dailyBill',
        name: '每日流水',
        component: DailyBill,
      },
      {
        path: '/tradingRecord',
        name: '交易记录',
        component: TradingRecord,
      },
    ],
  }

];

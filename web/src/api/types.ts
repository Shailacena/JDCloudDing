export interface IResponseBody<T> {
  success: boolean;
  code: number;
  message: string;
  data: T;
}

export interface Pagination {
  currentPage?: number;
  pageSize?: number;
}

export interface TableData<T> {
  list: Array<T>;
  total: number;
}

export interface AdminLoginReq {
  username: string;
  password: string;
  verifiCode: string;
}

export interface AdminLoginResp {
  id: number;
  token: string;
  nickname: string;
  role: number;
}

export interface AdminLogoutReq {}

export interface AdminLogoutResp {}

export interface AdminRegisterReq extends AdminBaseInfoReq {}

export interface AdminBaseInfoReq {
  username: string;
  nickname: string;
  remark: string;
  createdBy: number;
}

export interface AdminRegisterResp {
  username: string;
  nickname: string;
  password: string;
}

export interface ListAdminReq {}

export interface ListAdminResp extends TableData<IAdmin> {}

export interface IAdmin {
  id: number;
  username: string;
  nickname: string;
  remark: string;
  enable: number;
  role: number;
  urlKey: string;
  parentId: number;
  masterId: number;
}

export interface AdminResetPasswordReq {
  username: string;
}

export interface AdminResetPasswordResp {
  password: string;
}

export interface AdminSetPasswordReq {
  oldPassword: string;
  newPassword: string;
}

export interface AdminSetPasswordResp {}

export interface AdminResetVerifiCodeReq {
  id: number;
}

export interface AdminResetVerifiCodeResp {
  urlKey: string;
}

export interface AdminDeleteReq {
  username: string;
}

export interface AdminDeleteResp {}

export interface AdminUpdateReq extends AdminBaseInfoReq {}

export interface AdminUpdateResp {}

export interface AdminEnableReq {
  username: string;
  enable: number;
}

export interface AdminEnableResp {
  enable: number;
}

export interface ListOperationLogReq extends Pagination {}

export interface ListOperationLogResp extends TableData<IOperationLog> {}

export interface IOperationLog {
  ip: string;
  operation: string;
  operationStr: string;
  operator: number;
  operatorName: string;
  createAt: number;
}

export interface ListRealNameAccountReq extends Pagination {}

export interface ListRealNameAccountResp extends TableData<IRealNameAccount> {}

export interface IRealNameAccount extends BaseRealNameAccount {
  realNameCount: number;
  enable: number;
  remark: string;
}

export interface ListJDAccountReq extends JDAccountSearchParams, Pagination {}

export interface JDAccountSearchParams {
  id?: number;
  account?: string;
  status?: Array<number>;
  startAt?: string;
  endAt?: string;
}

export interface ListJDAccountResp extends TableData<IJDAccount> {}

export interface IJDAccount {
  id: number;
  account: string;
  realNameStatus: number;
  totalOrderCount: number;
  todayOrderCount: number;
  totalSuccessOrderCount: number;
  onlineStatus: number;
  status: number;
  remark: number;
  createAt: number;
}

export interface JDAccountDeleteReq extends JDAccountSearchParams {
  isAll: boolean;
}

export interface JDAccountDeleteResp {}

export interface JDAccountResetReq extends JDAccountSearchParams {}

export interface JDAccountResetResp {}

export interface ListPartnerReq extends Pagination {
  partnerId?: number;
  ignoreStatistics?: boolean;
}

export interface ListPartnerResp extends TableData<IPartner> {}

export interface IPartner {
  id: number;
  nickname: string;
  balance: number;
  priority: number;
  superiorAgent: string;
  level: number;
  stockAmount: number;
  enable: number;
  aqsAppSecret: string;
  aqsToken: string;
  remark: string;
  channelId: string;
  type: number;
  urlKey: string;
  parentId: number;
  darkNumberLength?: number;

  anssyAppSecret: string;
  anssyToken: string;
  anssyExpiredAt: number;

  todayOrderAmount: number;
  todayOrderNum: number;
  todaySuccessAmount: number;
  todaySuccessOrderNum: number;

  last1HourTotal: number;
  last1HourSuccess: number;
  last30MinutesTotal: number;
  last30MinutesSuccess: number;
}

export interface PartnerBaseInfoReq {
  nickname?: string;
  // superiorAgent?: string;
  type: number;
  level?: number;
  priority?: number;
  aqsAppSecret?: string;
  aqsToken?: string;
  anssyAppSecret?: string;
  anssyToken?: string;
  payType?: number;
  channelId?: string;
  remark?: string;
  darkNumberLength?: number;
}

export interface PartnerRegisterReq extends PartnerBaseInfoReq {}

export interface PartnerRegisterResp {
  nickname: string;
  password: string;
}

export interface PartnerLoginReq {
  username: string;
  password: string;
  verifiCode: string;
}

export interface PartnerLoginResp {
  id: number;
  token: string;
  level: number;
  nickname: string;
}

export interface PartnerLogoutReq {}

export interface PartnerLogoutResp {}

export interface PartnerSetPasswordReq {
  oldpassword: string;
  newpassword: string;
}

// ===== 服务器状态 / 合并 / 当日趋势 =====
export interface GetServerStatusReq {}

export interface GetServerStatusResp {
  dbVersion: string;
  orderTotal: number;
  orderToday: number;
  successTodayCount: number;
  successTodayAmount: number;
  partnerTotal: number;
  merchantTotal: number;
}

export interface ConsolidateReq {
  retentionDays?: number;
  batchSize?: number;
  dryRun?: boolean;
}

export interface GetTodayTrendReq {
  partnerId?: number;
  merchantId?: number;
}

export interface TrendPoint {
  time: string;
  orderCount: number;
  orderAmount: number;
  successCount: number;
  successAmount: number;
}

export interface GetTodayTrendResp {
  points: TrendPoint[];
}

export interface PartnerSetPasswordResp {}

export interface PartnerUpdateReq extends PartnerBaseInfoReq {
  id: number;
  nickname?: string;
  rechargeTime?: number;
  // stockAmount?: number;
  privateKey?: string;
  enable?: number;
}

export interface PartnerUpdateResp {}

export interface PartnerSyncGoodsReq {
  id: number;
}

export interface PartnerSyncGoodsResp {}

export interface PartnerResetPasswordReq {
  id: number;
}

export interface PartnerResetPasswordResp {
  password: string;
}

export interface PartnerResetVerifiCodeReq {
  id: number;
}

export interface PartnerResetVerifiCodeResp {
  urlKey: string;
}

export interface PartnerDeleteReq {
  id: number;
}

export interface PartnerDeleteResp {}

export interface ListPartnerBillResp {
  list: Array<IPartnerBill>;
}

export interface IPartnerBill {
  partnerId: number;
  type: number;
  changeMoney: number;
  money: number;
  remark: string;
  createAt: number;
}

export interface UpdatePartnerBalanceReq {
  adminId: number;
  partnerId: number;
  changeAmount: number;
  password: string;
}

export interface UpdatePartnerBalanceResp {}

export interface ListPartnerBalanceBillReq extends Pagination {
  partnerId?: number;
  startAt?: string;
  endAt?: string;
}

export interface ListPartnerBalanceBillResp
  extends TableData<IPartnerBalanceBill> {}

export interface IPartnerBalanceBill {
  id: number;
  partnerId: number;
  nickname: string;
  orderId: string;
  from: number;
  balance: number;
  changeAmount: number;
  createAt: number;
}

export interface ListGoodsReq extends Pagination {
  partnerId?: number;
}

export interface ListGoodsResp extends TableData<IGoods> {}

export interface IGoods extends GoodsInfo {
  id: number;
  createAt: number;
}

interface GoodsInfo {
  partnerId: number;
  skuId: string;
  // brandId: string;
  price: number;
  realPrice: number;
  shopName: string;
  status: number;
}

export interface GoodsCreateReq extends GoodsInfo {}

export interface GoodsCreateResp {}

export interface GoodsUpdateReq extends GoodsInfo {
  id: number;
}

export interface GoodsUpdateResp {}

export interface GoodsDeleteReq {
  id: number;
}

export interface GoodsDeleteResp {}

export interface ListMerchantReq extends Pagination {
  merchantId?: number;
  ignoreStatistics?: boolean;
}

export interface ListMerchantResp extends TableData<IMerchant> {}

export interface IMerchant {
  id: number;
  username: string;
  nickname: string;
  privateKey: string;
  createAt: number;
  balance: number;
  totalAmount: number;
  todayAmount: number;
  enable: number;
  remark: string;
  urlKey: string;
  parentId: number;
}

export interface MerchantRegisterReq {
  nickname: string;
  remark?: string;
}

export interface MerchantRegisterResp {
  nickname: string;
  password: string;
}

export interface MerchantResetPasswordReq {
  id: number;
}

export interface MerchantResetPasswordResp {
  password: string;
}

export interface AdminResetVerifiCodeReq {
  id: number;
}

export interface AdminResetVerifiCodeResp {
  urlKey: string;
}

export interface MerchantUpdateReq {
  username: string;
  nickname: string;
  remark?: string;
  isDel?: boolean;
}

export interface MerchantUpdateResp {}

export interface MerchantEnableReq {
  username: string;
  enable: number;
}

export interface MerchantEnableResp {
  enable: number;
}

export interface MerchantLoginReq {
  username: string;
  password: string;
  verifiCode: string;
}

export interface MerchantLoginResp {
  id: number;
  token: string;
  nickname: string;
}

export interface MerchantLogoutReq {}

export interface MerchantLogoutResp {}

export interface MerchantSetPasswordReq {
  oldPassword: string;
  newPassword: string;
}

export interface MerchantSetPasswordResp {}

export interface UpdateMerchantBalanceReq {
  adminId: number;
  merchantId: number;
  changeAmount: number;
  password: string;
}

export interface UpdateMerchantBalanceResp {}

export interface ListMerchantBalanceBillReq extends Pagination {
  merchantId?: number;
  startAt?: string;
  endAt?: string;
}

export interface ListMerchantBalanceBillResp
  extends TableData<IMerchantBalanceBill> {}

export interface IMerchantBalanceBill {
  id: number;
  merchantId: number;
  nickname: string;
  orderId: string;
  from: number;
  balance: number;
  changeAmount: number;
  createAt: number;
}

export interface MerchantGetBalanceReq {}

export interface MerchantGetBalanceResp {
  balance: number;
}

export interface ListDailyBillReq {
  partnerId?: number;
  merchantId?: number;
}

export interface ListDailyBillResp {
  list: Array<BaseDailyBill>;
}

export interface BaseDailyBill {
  date: string;
  totalOrderAmount: number;
  totalSuccessAmount: number;
  totalOrderNum: number;
  totalSuccessOrderNum: number;
}

export interface IDailyBill extends BaseDailyBill {
  id: number;
  nickname: string;
  balance: number;
  time: number;
}

export interface ListDailyBillByPartnerReq {
  partnerId?: number;
  startAt?: string;
  endAt?: string;
}

export interface ListDailyBillByPartnerResp extends TableData<IDailyBill> {}

export interface ListDailyBillByMerchantReq {
  merchantId?: number;
  startAt?: string;
  endAt?: string;
}

export interface ListDailyBillByMerchantResp extends TableData<IDailyBill> {}

export interface ListOrderReq extends Pagination {
  partnerId?: number;
  merchantId?: number;
  orderId?: string;
  partnerOrderId?: string;
  merchantOrderId?: string;
  startAt?: string;
  endAt?: string;
}

export interface ListOrderResp extends TableData<IOrder> {}

export interface IOrder {
  orderId: string;
  merchantOrderId: number;
  partnerOrderId: number;
  amount: number;
  receivedAmount: number;
  payType: number;
  channelId: string;
  payAccount: string;
  payAt: number;
  createAt: number;
  status: number;
  notifyStatus: number;
  merchantId: number;
  merchantName: string;
  partnerId: number;
  partnerName: string;
  skuId: string;
  shop: string;
  ip: string;
  remark: string;
}

export interface ConfirmOrderReq {
  orderId: string;
}

export interface ConfirmOrderResp {}

export interface GetOrderSummaryReq {}

export interface GetOrderSummaryResp {
  totalAmount?: number;
}

export interface ArchiveOrdersReq {
  adminId: number;
}

export interface ArchiveOrdersResp {
  archiveDate: string;
  totalAmount: number;
  orderCount: number;
}

export interface JDAccountCreateReq {
  accountList: Array<IJDAccountCreate>;
  remark: string;
}

export interface IJDAccountCreate {
  account: string;
  wsKey: string;
}

export interface JDAccountCreateResp {}

export interface JDAccountEnableReq {
  id: number;
  status: number;
}

export interface JDAccountEnableResp {}

export interface RealNameAccountCreateReq {
  accountList: Array<BaseRealNameAccount>;
  remark: string;
}

export interface BaseRealNameAccount {
  idNumber: string;
  name: string;
  mobile: string;
  address: string;
}

export interface RealNameAccountCreateResp {}

export interface GetMasterIncomeReq {
  masterId: number;
}

export interface GetMasterIncomeResp {
  totalIncome?: number;
}

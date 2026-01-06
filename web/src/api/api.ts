import { useAxios } from "./AxiosProvider";
import {
  AdminLoginReq,
  IResponseBody,
  AdminLoginResp,
  AdminLogoutReq,
  AdminLogoutResp,
  AdminRegisterReq,
  AdminRegisterResp,
  AdminResetPasswordReq,
  AdminResetPasswordResp,
  AdminResetVerifiCodeReq,
  AdminResetVerifiCodeResp,
  ListAdminResp,
  AdminDeleteReq,
  AdminDeleteResp,
  AdminUpdateReq,
  AdminUpdateResp,
  AdminEnableReq,
  AdminEnableResp,
  ListRealNameAccountResp,
  RealNameAccountCreateReq,
  RealNameAccountCreateResp,
  ListJDAccountResp,
  JDAccountCreateReq,
  JDAccountCreateResp,
  ListPartnerResp,
  PartnerRegisterReq,
  PartnerRegisterResp,
  PartnerLoginReq,
  PartnerLoginResp,
  PartnerSetPasswordReq,
  PartnerSetPasswordResp,
  PartnerResetPasswordReq,
  PartnerResetPasswordResp,
  PartnerDeleteReq,
  PartnerDeleteResp,
  PartnerUpdateReq,
  PartnerUpdateResp,
  ListGoodsResp,
  GoodsCreateReq,
  GoodsCreateResp,
  ListMerchantResp,
  MerchantRegisterReq,
  MerchantRegisterResp,
  MerchantLoginReq,
  MerchantLoginResp,
  MerchantSetPasswordReq,
  MerchantSetPasswordResp,
  ListDailyBillResp,
  ListOrderResp,
  AdminSetPasswordReq,
  AdminSetPasswordResp,
  JDAccountEnableReq,
  JDAccountEnableResp,
  ListJDAccountReq,
  JDAccountDeleteReq,
  JDAccountDeleteResp,
  JDAccountResetReq,
  JDAccountResetResp,
  GoodsDeleteReq,
  GoodsDeleteResp,
  GoodsUpdateReq,
  GoodsUpdateResp,
  MerchantUpdateReq,
  MerchantUpdateResp,
  ListOrderReq,
  PartnerSyncGoodsReq,
  PartnerSyncGoodsResp,
  MerchantEnableReq,
  MerchantEnableResp,
  MerchantLogoutReq,
  MerchantLogoutResp,
  ListMerchantBalanceBillReq,
  UpdateMerchantBalanceReq,
  UpdateMerchantBalanceResp,
  ListMerchantBalanceBillResp,
  UpdatePartnerBalanceReq,
  UpdatePartnerBalanceResp,
  ListPartnerBalanceBillReq,
  ListPartnerBalanceBillResp,
  ListGoodsReq,
  ListPartnerReq,
  ListMerchantReq,
  ConfirmOrderReq,
  ConfirmOrderResp,
  PartnerLogoutReq,
  PartnerLogoutResp,
  ListDailyBillReq,
  GetOrderSummaryReq,
  GetOrderSummaryResp,
  MerchantResetPasswordReq,
  MerchantResetPasswordResp,
  MerchantGetBalanceReq,
  MerchantGetBalanceResp,
  ListDailyBillByPartnerReq,
  ListDailyBillByPartnerResp,
  ListDailyBillByMerchantResp,
  ListDailyBillByMerchantReq,
  PartnerResetVerifiCodeReq,
  PartnerResetVerifiCodeResp,
  ListAdminReq,
  ListOperationLogReq,
  ListOperationLogResp,
  ListRealNameAccountReq,
  GetMasterIncomeReq,
  GetMasterIncomeResp,
  GetServerStatusReq,
  GetServerStatusResp,
  GetTodayTrendReq,
  GetTodayTrendResp,
  ArchiveOrdersReq,
  ArchiveOrdersResp,
} from "./types";

export function useApis() {
  const ax = useAxios();
  return {
    // 管理员
    adminLogin(data: AdminLoginReq): Promise<IResponseBody<AdminLoginResp>> {
      return ax.post("/admin/login", data).then((res) => {
        return res?.data;
      });
    },
    adminLogout(data: AdminLogoutReq): Promise<IResponseBody<AdminLogoutResp>> {
      return ax.post("/admin/logout", data).then((res) => {
        return res?.data;
      });
    },
    adminRegister(
      data: AdminRegisterReq
    ): Promise<IResponseBody<AdminRegisterResp>> {
      return ax.post("/admin/register", data).then((res) => {
        return res?.data;
      });
    },
    adminSetPassword(
      data: AdminSetPasswordReq
    ): Promise<IResponseBody<AdminSetPasswordResp>> {
      return ax.post("/admin/setPassword", data).then((res) => {
        return res?.data;
      });
    },
    adminResetPassword(
      data: AdminResetPasswordReq
    ): Promise<IResponseBody<AdminResetPasswordResp>> {
      return ax.post("/admin/resetPassword", data).then((res) => {
        return res?.data;
      });
    },
    adminResetVerifiCode(
      data: AdminResetVerifiCodeReq
    ): Promise<IResponseBody<AdminResetVerifiCodeResp>> {
      return ax.post("/admin/resetVerifiCode", data).then((res) => {
        return res?.data;
      });
    },
    listAdmin(params: ListAdminReq): Promise<IResponseBody<ListAdminResp>> {
      return ax.get("/admin/list", params).then((res) => {
        return res?.data;
      });
    },
    adminDelete(data: AdminDeleteReq): Promise<IResponseBody<AdminDeleteResp>> {
      return ax.post("/admin/delete", data).then((res) => {
        return res?.data;
      });
    },
    adminUpdate(data: AdminUpdateReq): Promise<IResponseBody<AdminUpdateResp>> {
      return ax.post("/admin/update", data).then((res) => {
        return res?.data;
      });
    },
    adminEnable(data: AdminEnableReq): Promise<IResponseBody<AdminEnableResp>> {
      return ax.post("/admin/enable", data).then((res) => {
        return res?.data;
      });
    },

    listOperationLog(params?: ListOperationLogReq): Promise<IResponseBody<ListOperationLogResp>> {
      return ax.get("/admin/operationLog", params).then((res) => {
        return res?.data;
      });
    },

    // 服务器状态与工具
    adminGetServerStatus(params?: GetServerStatusReq): Promise<IResponseBody<GetServerStatusResp>> {
      return ax.get("/admin/server/status", params).then((res) => {
        return res?.data;
      });
    },
    adminGetTodayTrend(params?: GetTodayTrendReq): Promise<IResponseBody<GetTodayTrendResp>> {
      return ax.get("/admin/order/trend/today", params).then((res) => {
        return res?.data;
      });
    },

    // 实名
    listRealNameAccount(
      params?: ListRealNameAccountReq
    ): Promise<IResponseBody<ListRealNameAccountResp>> {
      return ax.get("/realNameAccount/list", params).then((res) => {
        return res?.data;
      });
    },
    realNameAccountCreate(
      data: RealNameAccountCreateReq
    ): Promise<IResponseBody<RealNameAccountCreateResp>> {
      return ax.post("/realNameAccount/create", data).then((res) => {
        return res?.data;
      });
    },

    // 京东账号
    listJDAccount(
      params?: ListJDAccountReq
    ): Promise<IResponseBody<ListJDAccountResp>> {
      return ax.get("/jdAccount/list", params).then((res) => {
        return res?.data;
      });
    },
    jdAccountCreate(
      data: JDAccountCreateReq
    ): Promise<IResponseBody<JDAccountCreateResp>> {
      return ax.post("/jdAccount/create", data).then((res) => {
        return res?.data;
      });
    },
    jdAccountEnable(
      data: JDAccountEnableReq
    ): Promise<IResponseBody<JDAccountEnableResp>> {
      return ax.post("/jdAccount/enable", data).then((res) => {
        return res?.data;
      });
    },
    jdAccountDelete(
      data: JDAccountDeleteReq
    ): Promise<IResponseBody<JDAccountDeleteResp>> {
      return ax.post("/jdAccount/delete", data).then((res) => {
        return res?.data;
      });
    },
    jdAccountReset(
      data: JDAccountResetReq
    ): Promise<IResponseBody<JDAccountResetResp>> {
      return ax.post("/jdAccount/reset", data).then((res) => {
        return res?.data;
      });
    },

    // 合作商
    listPartner(
      data: ListPartnerReq
    ): Promise<IResponseBody<ListPartnerResp>> {
      return ax.get("/partner/list", data).then((res) => {
        return res?.data;
      });
    },
    partnerRegister(
      data: PartnerRegisterReq
    ): Promise<IResponseBody<PartnerRegisterResp>> {
      return ax.post("/partner/register", data).then((res) => {
        return res?.data;
      });
    },
    partnerResetPassword(
      data: PartnerResetPasswordReq
    ): Promise<IResponseBody<PartnerResetPasswordResp>> {
      return ax.post("/partner/resetPassword", data).then((res) => {
        return res?.data;
      });
    },
    partnerResetVerifiCode(
      data: PartnerResetVerifiCodeReq
    ): Promise<IResponseBody<PartnerResetVerifiCodeResp>> {
      return ax.post("/partner/resetVerifiCode", data).then((res) => {
        return res?.data;
      });
    },
    partnerDelete(
      data: PartnerDeleteReq
    ): Promise<IResponseBody<PartnerDeleteResp>> {
      return ax.post("/partner/delete", data).then((res) => {
        return res?.data;
      });
    },
    partnerUpdate(
      data: PartnerUpdateReq
    ): Promise<IResponseBody<PartnerUpdateResp>> {
      return ax.post("/partner/update", data).then((res) => {
        return res?.data;
      });
    },
    partnerUpdateBalance(
      data: UpdatePartnerBalanceReq
    ): Promise<IResponseBody<UpdatePartnerBalanceResp>> {
      return ax.post("/partner/updateBalance", data).then((res) => {
        return res?.data;
      });
    },
    partnerSyncGoods(
      data: PartnerSyncGoodsReq
    ): Promise<IResponseBody<PartnerSyncGoodsResp>> {
      return ax.post("/partner/syncGoods", data).then((res) => {
        return res?.data;
      });
    },
    listPartnerBalanceBill(
      params?: ListPartnerBalanceBillReq
    ): Promise<IResponseBody<ListPartnerBalanceBillResp>> {
      return ax.get("/partner/listBalanceBill", params).then((res) => {
        return res?.data;
      });
    },

    listGoods(data: ListGoodsReq): Promise<IResponseBody<ListGoodsResp>> {
      return ax.post("/goods/list", data).then((res) => {
        return res?.data;
      });
    },
    createGoods(data: GoodsCreateReq): Promise<IResponseBody<GoodsCreateResp>> {
      return ax.post("/goods/create", data).then((res) => {
        return res?.data;
      });
    },
    goodsUpdate(data: GoodsUpdateReq): Promise<IResponseBody<GoodsUpdateResp>> {
      return ax.post("/goods/update", data).then((res) => {
        return res?.data;
      });
    },
    goodsDelete(data: GoodsDeleteReq): Promise<IResponseBody<GoodsDeleteResp>> {
      return ax.post("/goods/delete", data).then((res) => {
        return res?.data;
      });
    },

    // 商户
    listMerchant(
      params: ListMerchantReq
    ): Promise<IResponseBody<ListMerchantResp>> {
      return ax.get("/merchant/list", params).then((res) => {
        return res?.data;
      });
    },
    merchantRegister(
      data: MerchantRegisterReq
    ): Promise<IResponseBody<MerchantRegisterResp>> {
      return ax.post("/merchant/register", data).then((res) => {
        return res?.data;
      });
    },
    merchantResetPassword(
      data: MerchantResetPasswordReq
    ): Promise<IResponseBody<MerchantResetPasswordResp>> {
      return ax.post("/merchant/resetPassword", data).then((res) => {
        return res?.data;
      });
    },
    merchantResetVerifiCode(
      data: AdminResetVerifiCodeReq
    ): Promise<IResponseBody<AdminResetVerifiCodeResp>> {
      return ax.post("/merchant/resetVerifiCode", data).then((res) => {
        return res?.data;
      });
    },
    merchantUpdate(
      data: MerchantUpdateReq
    ): Promise<IResponseBody<MerchantUpdateResp>> {
      return ax.post("/merchant/update", data).then((res) => {
        return res?.data;
      });
    },
    merchantUpdateBalance(
      data: UpdateMerchantBalanceReq
    ): Promise<IResponseBody<UpdateMerchantBalanceResp>> {
      return ax.post("/merchant/updateBalance", data).then((res) => {
        return res?.data;
      });
    },
    merchantEnable(data: MerchantEnableReq): Promise<IResponseBody<MerchantEnableResp>> {
      return ax.post("/merchant/enable", data).then((res) => {
        return res?.data;
      });
    },
    listMerchantBalanceBill(
      params?: ListMerchantBalanceBillReq
    ): Promise<IResponseBody<ListMerchantBalanceBillResp>> {
      return ax.get("/merchant/listBalanceBill", params).then((res) => {
        return res?.data;
      });
    },

    listDailyBill(
      params?: ListDailyBillReq
    ): Promise<IResponseBody<ListDailyBillResp>> {
      return ax.get("/statistics/listDailyBill", params).then((res) => {
        return res?.data;
      });
    },
    listDailyBillByPartner(
      params?: ListDailyBillByPartnerReq
    ): Promise<IResponseBody<ListDailyBillByPartnerResp>> {
      return ax.get("/statistics/listDailyBillByPartner", params).then((res) => {
        return res?.data;
      });
    },
    listDailyBillByMerchant(
      params?: ListDailyBillByMerchantReq
    ): Promise<IResponseBody<ListDailyBillByMerchantResp>> {
      return ax.get("/statistics/listDailyBillByMerchant", params).then((res) => {
        return res?.data;
      });
    },


    // 订单
    listOrder(
      params?: ListOrderReq
    ): Promise<IResponseBody<ListOrderResp>> {
      return ax.get("/order/list", params).then((res) => {
        return res?.data;
      });
    },

    confirmOrder(
      data: ConfirmOrderReq
    ): Promise<IResponseBody<ConfirmOrderResp>> {
      return ax.post("/order/confirm", data).then((res) => {
        return res?.data;
      });
    },
    getOrderSummary(
      params?: GetOrderSummaryReq
    ): Promise<IResponseBody<GetOrderSummaryResp>> {
      return ax.get("/order/summary", params).then((res) => {
        return res?.data;
      });
    },
    
    // 获取主账号总收入
    getMasterIncome(
      params: GetMasterIncomeReq
    ): Promise<IResponseBody<GetMasterIncomeResp>> {
      return ax.post("/admin/master/income", params).then((res) => {
          return res?.data;
        });
    },

    archiveOrders(
      data: ArchiveOrdersReq
    ): Promise<IResponseBody<ArchiveOrdersResp>> {
      return ax.post("/order/archive", data).then((res) => {
        return res?.data;
      });
    },

    // 合作商后台
    partner1Login(
      data: PartnerLoginReq
    ): Promise<IResponseBody<PartnerLoginResp>> {
      return ax.post("/partner1/login", data).then((res) => {
        return res?.data;
      });
    },
    partner1Logout(data: PartnerLogoutReq): Promise<IResponseBody<PartnerLogoutResp>> {
      return ax.post("/partner1/logout", data).then((res) => {
        return res?.data;
      });
    },
    partner1SetPassword(
      data: PartnerSetPasswordReq
    ): Promise<IResponseBody<PartnerSetPasswordResp>> {
      return ax.post("/partner1/setPassword", data).then((res) => {
        return res?.data;
      });
    },
    partner1SyncGoods(
      data: PartnerSyncGoodsReq
    ): Promise<IResponseBody<PartnerSyncGoodsResp>> {
      return ax.post("/partner1/syncGoods", data).then((res) => {
        return res?.data;
      });
    },
    // listPartner1Bill(): Promise<IResponseBody<ListPartnerBillResp>> {
    //   return ax.get("/partner1/listBill").then((res) => {
    //     return res?.data;
    //   });
    // },
    listPartner1BalanceBill(
      params?: ListPartnerBalanceBillReq
    ): Promise<IResponseBody<ListPartnerBalanceBillResp>> {
      return ax.get("/partner1/listBalanceBill", params).then((res) => {
        return res?.data;
      });
    },
    listPartner1Order(
      params?: ListOrderReq
    ): Promise<IResponseBody<ListOrderResp>> {
      return ax.get("/partner1/order/list", params).then((res) => {
        return res?.data;
      });
    },
    listPartner1StatisticsBill(
      params?: ListDailyBillReq
    ): Promise<IResponseBody<ListDailyBillResp>> {
      return ax.get("/partner1/statistics/listBill", params).then((res) => {
        return res?.data;
      });
    },
    listPartner1Goods(data: ListGoodsReq): Promise<IResponseBody<ListGoodsResp>> {
      return ax.post("/partner1/goods/list", data).then((res) => {
        return res?.data;
      });
    },

    // 商户后台
    merchant1Login(
      data: MerchantLoginReq
    ): Promise<IResponseBody<MerchantLoginResp>> {
      return ax.post("/merchant1/login", data).then((res) => {
        return res?.data;
      });
    },
    merchant1Logout(data: MerchantLogoutReq): Promise<IResponseBody<MerchantLogoutResp>> {
      return ax.post("/merchant1/logout", data).then((res) => {
        return res?.data;
      });
    },
    merchant1SetPassword(
      data: MerchantSetPasswordReq
    ): Promise<IResponseBody<MerchantSetPasswordResp>> {
      return ax.post("/merchant1/setPassword", data).then((res) => {
        return res?.data;
      });
    },
    listMerchant1BalanceBill(
      params?: ListMerchantBalanceBillReq
    ): Promise<IResponseBody<ListMerchantBalanceBillResp>> {
      return ax.get("/merchant1/listBalanceBill", params).then((res) => {
        return res?.data;
      });
    },
    listMerchant1StatisticsBill(
      params?: ListDailyBillReq
    ): Promise<IResponseBody<ListDailyBillResp>> {
      return ax.get("/merchant1/statistics/listBill", params).then((res) => {
        return res?.data;
      });
    },
    listMerchant1Order(
      params?: ListOrderReq
    ): Promise<IResponseBody<ListOrderResp>> {
      return ax.get("/merchant1/order/list", params).then((res) => {
        return res?.data;
      });
    },
    merchantGetBalance(
      params?: MerchantGetBalanceReq
    ): Promise<IResponseBody<MerchantGetBalanceResp>> {
      return ax.get("/merchant1/getBalance", params).then((res) => {
        return res?.data;
      });
    },
  };
}

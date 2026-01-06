import dayjs from "dayjs";
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';
import { EnableStatus } from "./constant";

// 扩展 dayjs 插件
dayjs.extend(utc);
dayjs.extend(timezone);

export enum OrderStatus {
  UnPaid = 1,
  Paid = 2,
  Finish = 3,
  RefundSuccessful = 4,
  RefundFailed = 5,
}

export function convertOrderStatus(status: number): string {
  switch (status) {
    case OrderStatus.Paid:
      return "已支付";
    case OrderStatus.Finish:
      return "已完成";
    case OrderStatus.RefundSuccessful:
      return "退款成功";
    case OrderStatus.RefundFailed:
      return "退款失败";

    default:
      return "待支付";
  }
}

export enum NotifyStatus {
  NotNotify = 1,
  Notify = 2,
}

export function convertNotifyStatus(status: number): string {
  switch (status) {
    case NotifyStatus.Notify:
      return "已通知";

    default:
      return "未通知";
  }
}

export function convertTimestamp(
  ts: number,
  format: string = "YYYY-MM-DD HH:mm:ss"
): string {
  if (ts > 0) {
    return dayjs.unix(ts).tz('Asia/Shanghai').format(format);
  }
  return "-";
}

export function convertEnable(enable: number): string {
  switch (enable) {
    case EnableStatus.Enabled:
      return "启用";

    case EnableStatus.Disabled:
      return "禁用";

    default:
      return "";
  }
}

export enum PayType {
  AppPay = 1,
  WebPay = 2,
  AliPay = 3,
}

export interface IPayType {
  payType: PayType;
  label: string;
}

export let AllPayType: IPayType[] = [
  {
    payType: PayType.AppPay,
    label: "App跳转",
  },
  {
    payType: PayType.WebPay,
    label: "Web跳转",
  },
  {
    payType: PayType.AliPay,
    label: "支付宝跳转",
  },
];

export function convertPayType(payType: PayType): string {
  for (let index = 0; index < AllPayType.length; index++) {
    const p = AllPayType[index];
    if (p.payType === payType) {
      return p.label;
    }
  }
  return "";
}

export function isJDShop(channelId?: string): boolean {
  switch (channelId) {
    case ChannelId.ChannelJDPay:
      return true;

    default:
      return false;
  }
}

export enum ChannelId {
  ChannelJDGame = "JD000000", // 京东游戏
  ChannelJDEntity = "JD111111", // 京东实物
  ChannelTBShop = "TB000000", // 淘宝店铺
  ChannelTBQRCode = "TB111111", // 淘宝码上收
  ChannelTBClipboard = "TB222222", // 淘宝复制
  ChannelTBPayForAnother = "TB668888", // 淘宝代付
  ChannelTBDirectPay = "TB686088", // 天猫直付
  ChannelTBECoupon = "TB087888", // 淘宝电子券
  ChannelJDPay = "JS000000", // 京东复制
  ChannelJDCK = "JS111111", // 京东ck
}

export interface IChannelId {
  channelId: ChannelId;
  label: string;
}

interface IChannelPayType {
  [key: string]: any;
}

export let ChannelPayType: IChannelPayType = {};
ChannelPayType[ChannelId.ChannelTBDirectPay] = [
  {
    payType: PayType.AppPay,
    label: "App跳转",
  },
  {
    payType: PayType.WebPay,
    label: "Web跳转",
  },
  {
    payType: PayType.AliPay,
    label: "支付宝跳转",
  },
];

ChannelPayType[ChannelId.ChannelJDPay] = [
  {
    payType: PayType.AppPay,
    label: "App跳转",
  },
];

ChannelPayType[ChannelId.ChannelJDCK] = [
  {
    payType: PayType.AppPay,
    label: "App跳转",
  },
];

export let AllChannelId: IChannelId[] = [
  // {
  //   channelId: ChannelId.ChannelTBClipboard,
  //   label: "淘宝复制",
  // },
  // {
  //   channelId: ChannelId.ChannelTBPayForAnother,
  //   label: "淘宝代付",
  // },
  {
    channelId: ChannelId.ChannelTBDirectPay,
    label: "天猫直付",
  },
  // {
  //   channelId: ChannelId.ChannelTBECoupon,
  //   label: "淘宝电子券",
  // },
  {
    channelId: ChannelId.ChannelJDPay,
    label: "京东复制",
  },
  {
    channelId: ChannelId.ChannelJDCK,
    label: "京东ck",
  },
  // {
  //   channelId: ChannelId.ChannelJDGame,
  //   label: "京东游戏",
  // },
  // {
  //   channelId: ChannelId.ChannelJDEntity,
  //   label: "京东实物",
  // },
  // {
  //   channelId: ChannelId.ChannelTBShop,
  //   label: "淘宝店铺",
  // },
  // {
  //   channelId: ChannelId.ChannelTBQRCode,
  //   label: "淘宝码上收",
  // },
];

export function convertChannelId(channelId: string): string {
  for (let index = 0; index < AllChannelId.length; index++) {
    const p = AllChannelId[index];
    if (p.channelId === channelId) {
      return p.label;
    }
  }
  return "";
}

export enum BalanceFromType {
  BalanceFromTypeOrderAdd = 1,
  BalanceFromTypeOrderDeduct = 2,
  BalanceFromTypeSystemAdd = 3,
  BalanceFromTypeSystemDeduct = 4,
}

export function convertBalanceFrom(from: BalanceFromType): string {
  switch (from) {
    case BalanceFromType.BalanceFromTypeOrderAdd:
      return "订单收入";
    case BalanceFromType.BalanceFromTypeOrderDeduct:
      return "订单扣减";
    case BalanceFromType.BalanceFromTypeSystemAdd:
      return "平台增加";
    case BalanceFromType.BalanceFromTypeSystemDeduct:
      return "平台扣减";
    default:
      return "";
  }
}

export enum GoodsStatus {
  Enabled = 1,
  Disabled = 2,
}

export interface IPartnerType {
  id: number;
  label: string;
}

export enum PartnerType {
  Agiso = 1,
  Anssy = 2,
}

export let AllPartnerType: IPartnerType[] = [
  {
    id: PartnerType.Agiso,
    label: "阿奇索",
  },
  {
    id: PartnerType.Anssy,
    label: "安式",
  },
];

export enum JDAccountStatus {
  Normal = 1,
  Invalid = 2,
  Hot = 3,
  AddAddressErr = 4,
  SubmitOrderErr = 5,
  GetWxPayErr = 6,
}

export function convertJDAccountStatus(from: JDAccountStatus): string {
  switch (from) {
    case JDAccountStatus.Normal:
      return "启动";
    case JDAccountStatus.Invalid:
      return "过期";
    case JDAccountStatus.Hot:
      return "火热";
    case JDAccountStatus.AddAddressErr:
      return "添加地址失败";
    case JDAccountStatus.SubmitOrderErr:
      return "下单失败";
    case JDAccountStatus.GetWxPayErr:
      return "提码失败";
    default:
      return "";
  }
}

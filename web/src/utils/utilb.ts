import { EnableStatus } from "./constant"

export const isEnable = (enable: number): boolean => {
  return enable === EnableStatus.Enabled
}

export function toPercent(num: number): number {
  return Math.round(num * 100)
}

export function toPercentWithFixed(num1: number, num2: number): number {
  if (num2 === 0) {
    return 0
  }
  let num: number = +((num1 / num2).toFixed(4))
  return +(num * 100).toFixed(2)
}

export function formatNumberWithCommasAndDecimals(num: number) {
  return Number(num).toFixed(2).replace(/\B(?=(\d{3})+(?!\d))/g, ',');
}
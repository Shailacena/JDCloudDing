import { useCallback, useEffect, useState } from 'react';
import { Card, Table } from 'antd';
import type { TableProps } from 'antd';
import { useApis } from '../api/api';
import { IDailyBill, ListDailyBillByMerchantReq } from '../api/types';
import { formatNumberWithCommasAndDecimals, toPercentWithFixed } from '../utils/utilb';
import MerchantDailyBillSearchForm from './searchform/MerchantDailyBillSearchForm';
import { convertTimestamp } from '../utils/biz';
import { Dayjs } from 'dayjs';

interface DataType extends IDailyBill {
  key: number;
}

export interface IMerchantSearchResult {
  merchantId?: number;
  startAt?: Dayjs;
  endAt?: Dayjs;
}

const columns: TableProps<DataType>['columns'] = [
  {
    title: '日期', dataIndex: 'time', key: 'time', align: 'center', render: (_, d) => (
      convertTimestamp(d.time, 'YYYY-MM-DD')
    )
  },
  {
    title: '商户号', dataIndex: 'id', key: 'id', align: 'center',
  },
  {
    title: '商户名称', dataIndex: 'nickname', key: 'nickname', align: 'center',
  },
  {
    title: '余额', key: 'balance', dataIndex: 'balance', align: 'center', render: (_, d) => (
      formatNumberWithCommasAndDecimals(d.balance)
    )
  },
  {
    title: '订单总额', dataIndex: 'totalOrderAmount', key: 'totalOrderAmount', align: 'center', render: (_, d) => (
      formatNumberWithCommasAndDecimals(d.totalOrderAmount)
    )
  },
  {
    title: '成功订单总额', dataIndex: 'totalSuccessAmount', key: 'totalSuccessAmount', align: 'center', render: (_, d) => (
      formatNumberWithCommasAndDecimals(d.totalSuccessAmount)
    )
  },
  {
    title: '订单总数', key: 'totalOrderNum', dataIndex: 'totalOrderNum', align: 'center',
  },
  {
    title: '成功订单总数', key: 'totalSuccessOrderNum', dataIndex: 'totalSuccessOrderNum', align: 'center',
  },
  {
    title: '订单成功率', align: 'center', render: (_, d) => (
      toPercentWithFixed(d.totalSuccessOrderNum, d.totalOrderNum) + "%"
    )
  },
  {
    title: '订单金额成功率', align: 'center', render: (_, d) => (
      toPercentWithFixed(d.totalSuccessAmount, d.totalOrderAmount) + "%"
    )
  },
];

function MerchantDailyBill() {
  const [list, setList] = useState<DataType[]>([])
  let { listDailyBillByMerchant } = useApis()
  const [listLoadingStates, setListLoadingStates] = useState(false);
  const [reqParams, setReqParams] = useState<ListDailyBillByMerchantReq>();

  const fetchListMerchantBill = useCallback(async () => {
    try {
      setListLoadingStates(true)

      const { data } = await listDailyBillByMerchant(reqParams)
      let d: DataType[] = data?.list?.map((item, index) => {
        let newItem: DataType = {
          key: index,
          ...item
        }
        return newItem
      })

      setList(d)
    } catch (e) {
      console.error(e);
    } finally {
      setListLoadingStates(false)
    }
  }, [reqParams])

  useEffect(() => {
    fetchListMerchantBill()
  }, [reqParams])

  const onSearch = (value: ListDailyBillByMerchantReq) => {
    setReqParams({ ...value })
  }

  return (
    <>
      <Card>
        <div style={{ display: 'Flex' }}>
          <MerchantDailyBillSearchForm onSearch={onSearch} />
        </div>
        <Table<DataType>
          bordered
          size='middle'
          columns={columns}
          dataSource={list}
          scroll={{ x: 'max-content' }}
          loading={listLoadingStates}
          pagination={false}
        />
      </Card>
    </>
  )
}

export default MerchantDailyBill
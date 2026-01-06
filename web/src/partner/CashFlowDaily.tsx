import { Table, Card, Button } from 'antd';
import type { TableProps } from 'antd';
import { BaseDailyBill, ListDailyBillReq } from '../api/types';
import { useEffect, useState } from 'react';
import { useApis } from '../api/api';
import { useAppContext } from '../AppProvider';
import { formatNumberWithCommasAndDecimals } from '../utils/utilb';

interface DataType extends BaseDailyBill {
  key: string;
}

const columns: TableProps<DataType>['columns'] = [
  {
    title: '日期', dataIndex: 'date', key: 'date', align: 'center',
  },
  {
    title: '订单总金额', dataIndex: 'totalOrderAmount', key: 'totalOrderAmount', align: 'center', render: (_, d) => (
      formatNumberWithCommasAndDecimals(d.totalOrderAmount)
    )
  },
  {
    title: '成功总金额', dataIndex: 'totalSuccessAmount', key: 'totalSuccessAmount', align: 'center', render: (_, d) => (
      formatNumberWithCommasAndDecimals(d.totalSuccessAmount)
    )
  },
  {
    title: '总订单数', key: 'totalOrderNum', dataIndex: 'totalOrderNum', align: 'center',
  },
];

function CashFlowDaily() {
  const [list, setList] = useState<DataType[]>([])
  let { listPartner1StatisticsBill } = useApis()
  const [listLoadingStates, setListLoadingStates] = useState(false);
  const ctx = useAppContext();

  const fetchListStatisticsBill = async () => {
    try {
      setListLoadingStates(true)

      let params: ListDailyBillReq = {
        partnerId: ctx.cookie.id,
      }

      const { data } = await listPartner1StatisticsBill(params)
      let d: DataType[] = data?.list?.map((item, index) => {
        let newItem: DataType = {
          key: index.toString(),
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
  }

  useEffect(() => {
    fetchListStatisticsBill()
  }, [])

  return (
    <>
      <Card>
        <div style={{ marginBottom: '10px', display: 'Flex' }}>
          <Button type="primary" onClick={fetchListStatisticsBill} >刷新</Button>
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

export default CashFlowDaily
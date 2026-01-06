import { useEffect, useState } from 'react';
import { Button, Card, Table } from 'antd';
import type { TableProps } from 'antd';
import { useApis } from '../api/api';
import { BaseDailyBill } from '../api/types';
import { formatNumberWithCommasAndDecimals, toPercentWithFixed } from '../utils/utilb';


interface DataType extends BaseDailyBill {
  key: string;
}

const columns: TableProps<DataType>['columns'] = [
  {
    title: '日期', dataIndex: 'date', key: 'date', align: 'center',
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

function DailyBill() {
  const [list, setList] = useState<DataType[]>([])
  const [listLoadingStates, setListLoadingStates] = useState(false);
  let { listDailyBill: listStatisticsBill } = useApis()

  const fetchListStatisticsBill = async () => {
    try {
      setListLoadingStates(true)

      const { data } = await listStatisticsBill()
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

export default DailyBill
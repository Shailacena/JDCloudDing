import { useCallback, useEffect, useState } from 'react';
import { Table, Card, Typography, Flex } from 'antd';
import type { TableProps } from 'antd';
import { useApis } from '../api/api';
import { IMerchantBalanceBill, ListMerchantBalanceBillReq } from '../api/types';
import { convertBalanceFrom, convertTimestamp } from '../utils/biz';
import { PAGE_DEFAULT_INDEX, PAGE_SIZE } from '../components/types';
import BalanceBillSearchForm, { ListBalanceBill } from './searchform/BalanceBillSearchForm ';
import { formatNumberWithCommasAndDecimals } from '../utils/utilb';
import { useAppContext } from '../AppProvider';

interface DataType extends IMerchantBalanceBill {
  key: number;
}

export function convertChangeAmount(changeAmount: number) {
  if (changeAmount > 0) {
    return (
      <Typography.Text type="danger">+{formatNumberWithCommasAndDecimals(changeAmount)}</Typography.Text>
    )
  } else if (changeAmount < 0) {
    return (
      <Typography.Text type="success">{formatNumberWithCommasAndDecimals(changeAmount)}</Typography.Text>
    )
  } else {
    return <span>{formatNumberWithCommasAndDecimals(changeAmount)}</span>
  }
}

function BalanceBill() {
  const columns: TableProps<DataType>['columns'] = [
    {
      title: '订单号', key: 'orderId', dataIndex: 'orderId', align: 'center', render: (_, d) => (
        d.orderId || '-'
      )
    },
    {
      title: '类型', key: 'from', dataIndex: 'from', align: 'center', render: (_, d) => (
        convertBalanceFrom(d.from)
      )
    },
    {
      title: '账户余额', key: 'balance', dataIndex: 'balance', align: 'center', render: (_, d) => (
        formatNumberWithCommasAndDecimals(d.balance)
      )
    },
    {
      title: '交易金额', key: 'changeAmount', dataIndex: 'changeAmount', align: 'center', render: (_, d) => (
        convertChangeAmount(d.changeAmount)
      )
    },
    {
      title: '交易时间', key: 'createAt', dataIndex: 'createAt', align: 'center', render: (_, d) => (
        convertTimestamp(d.createAt)
      )
    },
  ];
  const [list, setList] = useState<DataType[]>([])
  let { listMerchant1BalanceBill } = useApis()
  const [total, setTotal] = useState(0);
  const [listLoadingStates, setListLoadingStates] = useState(false);
  const ctx = useAppContext()
  const [reqParams, setReqParams] = useState<ListMerchantBalanceBillReq>({
    currentPage: PAGE_DEFAULT_INDEX,
    pageSize: PAGE_SIZE.TEN,
    merchantId: ctx.cookie.id
  });

  const fetchMerchantBalanceBillList = useCallback(async () => {
    try {
      setListLoadingStates(true)

      const { data } = await listMerchant1BalanceBill(reqParams)
      let d: DataType[] = data?.list?.map((item, index) => {
        let newItem: DataType = {
          key: index,
          ...item
        }
        return newItem
      })
      setList(d)
      setTotal(data.total)

    } catch (e) {
      console.error(e);
    } finally {
      setListLoadingStates(false)
    }
  }, [reqParams])

  useEffect(() => {
    fetchMerchantBalanceBillList()
  }, [reqParams])

  const onSearch = (value: ListBalanceBill) => {
    setReqParams({ ...reqParams, ...value, currentPage: 1, pageSize: reqParams.pageSize })
  }

  const handleTableChange = (current: number, pageSize: number) => {
    setReqParams({ ...reqParams, currentPage: current, pageSize })
  };

  return (
    <>
      <Card>
        <Flex vertical>
          <BalanceBillSearchForm OnSearch={onSearch} />
        </Flex>

        <Table<DataType>
          bordered
          size='middle'
          pagination={{
            current: reqParams.currentPage,
            pageSize: reqParams.pageSize,
            total: total,
            onChange: handleTableChange,
          }}
          columns={columns}
          dataSource={list}
          scroll={{ x: 'max-content' }}
          loading={listLoadingStates}
        />
      </Card>
    </>
  )
}

export default BalanceBill
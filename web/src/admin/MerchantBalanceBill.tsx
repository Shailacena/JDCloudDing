import { useCallback, useEffect, useState } from 'react';
import { Table, Card, Typography, Flex } from 'antd';
import type { TableProps } from 'antd';
import { useApis } from '../api/api';
import { IMerchant, IMerchantBalanceBill, ListMerchantBalanceBillReq } from '../api/types';
import { convertBalanceFrom, convertTimestamp } from '../utils/biz';
import BalanceBillSearchForm, { ListBalanceBill, SearchFormIdItem, SearchType } from './searchform/BalanceBillSearchForm ';
import { PAGE_DEFAULT_INDEX, PAGE_SIZE } from '../components/types';
import { formatNumberWithCommasAndDecimals } from '../utils/utilb';

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

function MerchantBalanceBill() {
  const columns: TableProps<DataType>['columns'] = [
    {
      title: '商户ID', dataIndex: 'merchantId', key: 'merchantId', align: 'center',
    },
    {
      title: '商户名称', dataIndex: 'nickname', key: 'nickname', align: 'center',
    },
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
  const [balanceBillSearchIds, setBalanceBillSearchIds] = useState<SearchFormIdItem[]>([])
  let { listMerchantBalanceBill, listMerchant } = useApis()
  const [total, setTotal] = useState(0);
  const [listLoadingStates, setListLoadingStates] = useState(false);

  const [reqParams, setReqParams] = useState<ListMerchantBalanceBillReq>({
    currentPage: PAGE_DEFAULT_INDEX,
    pageSize: PAGE_SIZE.TEN
  });

  const fetchMerchantBalanceBillList = useCallback(async () => {
    try {
      setListLoadingStates(true)

      const { data } = await listMerchantBalanceBill(reqParams)
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

  const fetchListMerchant = async () => {
    const { data } = await listMerchant({ignoreStatistics: true})
    setMerchantIds(data.list)
  }

  const setMerchantIds = (datas: IMerchant[]) => {
    let ids: SearchFormIdItem[] = datas?.map((item) => {
      let newItem: SearchFormIdItem = {
        id: item.id,
        label: item.username + '(' + item.nickname + ')',
      }
      return newItem
    })
    setBalanceBillSearchIds(ids)
  }

  useEffect(() => {
    fetchListMerchant()
  }, [])

  useEffect(() => {
    fetchMerchantBalanceBillList()
  }, [reqParams])

  const onSearch = (value: ListBalanceBill) => {
    let { id, ...params } = value
    setReqParams({ ...params, merchantId: id, currentPage: 1, pageSize: reqParams.pageSize })
  }

  const handleTableChange = (current: number, pageSize: number) => {
    setReqParams({ ...reqParams, currentPage: current, pageSize })
  };

  return (
    <>
      <Card>
        <Flex vertical>
          <BalanceBillSearchForm ids={balanceBillSearchIds} OnSearch={onSearch} searchType={SearchType.Merchant} />
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

export default MerchantBalanceBill
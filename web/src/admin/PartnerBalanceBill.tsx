import { useCallback, useEffect, useState } from 'react';
import { Table, Card, Flex } from 'antd';
import type { TableProps } from 'antd';
import { useApis } from '../api/api';
import { IPartner, IPartnerBalanceBill, ListPartnerBalanceBillReq } from '../api/types';
import { convertBalanceFrom, convertTimestamp } from '../utils/biz';
import { convertChangeAmount } from './MerchantBalanceBill';
import BalanceBillSearchForm, { ListBalanceBill, SearchFormIdItem, SearchType } from './searchform/BalanceBillSearchForm ';
import { useAppContext } from '../AppProvider';
import { PAGE_DEFAULT_INDEX, PAGE_SIZE } from '../components/types';
import { formatNumberWithCommasAndDecimals } from '../utils/utilb';

interface DataType extends IPartnerBalanceBill {
  key: number;
}

function PartnerBalanceBill() {
  const columns: TableProps<DataType>['columns'] = [
    {
      title: '合作商ID', dataIndex: 'partnerId', key: 'partnerId', align: 'center',
    },
    {
      title: '合作商名称', dataIndex: 'nickname', key: 'nickname', align: 'center',
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
  let { listPartnerBalanceBill, listPartner } = useApis()
  const [total, setTotal] = useState(0);
  const ctx = useAppContext();
  const [listLoadingStates, setListLoadingStates] = useState(false);
  const [reqParams, setReqParams] = useState<ListPartnerBalanceBillReq>({
    currentPage: PAGE_DEFAULT_INDEX,
    pageSize: PAGE_SIZE.TEN
  });

  const fetchPartnerBalanceBillList = useCallback(async () => {
    try {
      setListLoadingStates(true)

      const { data } = await listPartnerBalanceBill(reqParams)
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
    fetchPartnerBalanceBillList()
  }, [reqParams])

  useEffect(() => {
    fetchListPartner()
  }, [])

  const fetchListPartner = async () => {
    const { data } = await listPartner({ ignoreStatistics: true })
    ctx.partnerList = data.list;
    setPartnerIds(data.list)
  }

  const setPartnerIds = (datas: IPartner[]) => {
    let ids: SearchFormIdItem[] = datas?.map((item) => {
      let newItem: SearchFormIdItem = {
        id: item.id,
        label: item.id + '(' + item.nickname + ')',
      }
      return newItem
    })
    setBalanceBillSearchIds(ids)
  }

  const handleTableChange = (current: number, pageSize: number) => {
    setReqParams({ ...reqParams, currentPage: current, pageSize })
  };

  const onSearch = (value: ListBalanceBill) => {
    let { id, ...params } = value
    setReqParams({ ...params, partnerId: id, currentPage: 1, pageSize: reqParams.pageSize })
  }

  return (
    <>
      <Card>
        <Flex vertical>
          <BalanceBillSearchForm ids={balanceBillSearchIds} OnSearch={onSearch} searchType={SearchType.Partner} />
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

export default PartnerBalanceBill
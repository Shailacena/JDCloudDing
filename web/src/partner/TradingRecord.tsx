import { useCallback, useEffect, useState } from 'react';
import { Card, Table } from 'antd';
import type { TableProps } from 'antd';
import { useApis } from '../api/api';
import { convertNotifyStatus, convertOrderStatus, convertPayType, convertTimestamp } from '../utils/biz';
import { IOrder, ListOrderReq } from '../api/types';
import { PAGE_DEFAULT_INDEX, PAGE_SIZE } from '../components/types';
import TradingRecordSearchForm from './searchform/TradingRecordSearchForm';
import { formatNumberWithCommasAndDecimals } from '../utils/utilb';
import { useAppContext } from '../AppProvider';


interface DataType extends IOrder {
  key: string;
}

function TradingRecord() {
  const columns: TableProps<DataType>['columns'] = [
    {
      title: '店铺订单号', dataIndex: 'partnerOrderId', key: 'partnerOrderId', align: 'center', render: (_, d) => {
        return d.partnerOrderId || '-'
      },
    },
    {
      title: '系统订单号', dataIndex: 'orderId', key: 'orderId', align: 'center',
    },
    {
      title: '合作商名称', key: 'shop', dataIndex: 'shop', align: 'center',
    },
    {
      title: '订单金额', key: 'amount', dataIndex: 'amount', align: 'center', render: (_, d) => (
        formatNumberWithCommasAndDecimals(d.amount)
      )
    },
    {
      title: '店铺名称', key: 'shop', dataIndex: 'shop', align: 'center',
    },
    {
      title: 'skuId', key: 'skuId', dataIndex: 'skuId', align: 'center',
    },
    {
      title: '订单状态', key: 'status', dataIndex: 'status', render: (_, d) => { return convertOrderStatus(d.status) }, align: 'center',
    },
    {
      title: '通知状态', key: 'notifyStatus', dataIndex: 'notifyStatus', render: (_, d) => {
        return convertNotifyStatus(d.notifyStatus)
      },
      align: 'center',
    },
    {
      title: '下单时间', key: 'createAt', dataIndex: 'createAt', render: (_, d) => {
        return convertTimestamp(d.createAt)
      },
      align: 'center',
    },
    {
      title: '下单账号', key: 'payAccount', dataIndex: 'payAccount', align: 'center', render: (_, d) => {
        return d.payAccount || '-'
      },
    },
    {
      title: '支付类型', key: 'payType', dataIndex: 'payType', render: (_, d) => {
        return convertPayType(d.payType)
      },
      align: 'center',
    },
  ];
  const [list, setList] = useState<DataType[]>([])
  let { listPartner1Order } = useApis()
  const [total, setTotal] = useState(0);
  const [listLoadingStates, setListLoadingStates] = useState(false);
  let ctx = useAppContext();
  const [reqParams, setReqParams] = useState<ListOrderReq>({
    currentPage: PAGE_DEFAULT_INDEX,
    pageSize: PAGE_SIZE.TEN,
    partnerId: ctx.cookie.id
  });

  const fetchListOrder = useCallback(async () => {
    try {
      setListLoadingStates(true)

      const { data } = await listPartner1Order(reqParams)
      let d: DataType[] = data?.list?.map((item, index) => {
        let newItem: DataType = {
          key: index.toString(),
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

  const handleTableChange = (current: number, pageSize: number) => {
    setReqParams({ ...reqParams, currentPage: current, pageSize })
  };

  const onSearch = (value: ListOrderReq) => {
    setReqParams({ ...reqParams, ...value, currentPage: 1, pageSize: reqParams.pageSize })
  }

  useEffect(() => {
    fetchListOrder()
  }, [reqParams])

  return (
    <>
      <Card>
        <div style={{ display: 'Flex' }}>
          <TradingRecordSearchForm OnSearch={onSearch} />
        </div>
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

export default TradingRecord
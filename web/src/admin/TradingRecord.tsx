import { useCallback, useEffect, useState } from 'react';
import { Button, Card, message, Popconfirm, Space, Table } from 'antd';
import type { TableProps } from 'antd';
import { useApis } from '../api/api';
import TradingRecordSearchForm from './searchform/TradingRecordSearchForm';
import { convertChannelId, convertNotifyStatus, convertOrderStatus, convertPayType, convertTimestamp, NotifyStatus, OrderStatus } from '../utils/biz';
import { IOrder, ListOrderReq } from '../api/types';
import { PAGE_DEFAULT_INDEX, PAGE_SIZE } from '../components/types';
import { formatNumberWithCommasAndDecimals } from '../utils/utilb';

interface DataType extends IOrder {
  key: string;
}

function TradingRecord() {
  const columns: TableProps<DataType>['columns'] = [
    {
      title: '系统订单号', dataIndex: 'orderId', key: 'orderId', align: 'center',
    },
    {
      title: '店铺订单号', dataIndex: 'partnerOrderId', key: 'partnerOrderId', align: 'center', render: (_, d) => {
        return d.partnerOrderId || '-'
      },
    },
    {
      title: '商户订单号', dataIndex: 'merchantOrderId', key: 'merchantOrderId', align: 'center',
    },
    {
      title: '商户号', dataIndex: 'merchantId', key: 'merchantId', align: 'center',
    },
    {
      title: '商户名称', dataIndex: 'merchantName', key: 'merchantName', align: 'center',
    },
    {
      title: '订单金额', key: 'amount', dataIndex: 'amount', align: 'center', render: (_, d) => (
        formatNumberWithCommasAndDecimals(d.amount)
      )
    },
    {
      title: '实收金额', key: 'receivedAmount', dataIndex: 'receivedAmount', align: 'center', render: (_, d) => (
        formatNumberWithCommasAndDecimals(d.receivedAmount)
      )
    },
    {
      title: '订单状态', key: 'status', dataIndex: 'status', align: 'center', render: (_, d) => {
        return convertOrderStatus(d.status)
      },

    },
    {
      title: '通知状态', key: 'notifyStatus', dataIndex: 'notifyStatus', align: 'center', render: (_, d) => {
        return convertNotifyStatus(d.notifyStatus)
      },
    },
    {
      title: '下单时间', key: 'createAt', dataIndex: 'createAt', align: 'center', render: (_, d) => {
        return convertTimestamp(d.createAt)
      },
    },
    {
      title: '支付类型', key: 'payType', dataIndex: 'payType', align: 'center', render: (_, d) => {
        return convertPayType(d.payType)
      },
    },
    {
      title: '支付通道', key: 'channelId', dataIndex: 'channelId', align: 'center', render: (_, d) => {
        return convertChannelId(d.channelId)
      },
    },
    {
      title: '下单账号', key: 'payAccount', dataIndex: 'payAccount', align: 'center', render: (_, d) => {
        return d.payAccount || '-'
      },
    },
    {
      title: 'skuId', key: 'skuId', dataIndex: 'skuId', align: 'center',
    },
    // {
    //   title: '店铺名称',  //   key: 'shop',  //   dataIndex: 'shop',  //   align: 'center',  // },
    {
      title: '合作商号', dataIndex: 'partnerId', key: 'partnerId', align: 'center',
    },
    {
      title: '合作商名称', key: 'partnerName', dataIndex: 'partnerName', align: 'center',
    },
    {
      title: 'ip', key: 'ip', dataIndex: 'ip', align: 'center',
    },
    {
      title: '备注', key: 'remark', dataIndex: 'remark', align: 'center', render: (_, d) => {
        return d.remark || '-'
      },
    },
    {
      title: '操作', key: 'action', align: 'center', fixed: 'right', render: (_, d) => (
        <Space size="middle">
          <Popconfirm title="警告" description="请确认收货？" onConfirm={() => { onConfirm(d.orderId) }} >
            <Button disabled={d.status === OrderStatus.Finish && d.notifyStatus == NotifyStatus.Notify} type="primary" size='small'>确认收货</Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];
  const [list, setList] = useState<DataType[]>([])
  let { listOrder, confirmOrder } = useApis()
  const [total, setTotal] = useState(0);
  const [listLoadingStates, setListLoadingStates] = useState(false);
  const [reqParams, setReqParams] = useState<ListOrderReq>({
    currentPage: PAGE_DEFAULT_INDEX,
    pageSize: PAGE_SIZE.TEN
  });

  const fetchListOrder = useCallback(async () => {
    try {
      setListLoadingStates(true)

      const { data } = await listOrder(reqParams)
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
    setReqParams({ ...value, currentPage: 1, pageSize: reqParams.pageSize })
  }

  useEffect(() => {
    fetchListOrder()
  }, [reqParams])

  const onConfirm = async (orderId: string) => {
    try {
      await confirmOrder({ orderId })

      message.success("确认收货成功")
      fetchListOrder()
    } catch (e) {
      console.error(e);
    }
  }

  // 纯JavaScript实现CSV导出
  const convertToCSV = (data: any[], headers: string[]) => {
    // CSV头部
    const csvHeaders = headers.join(',');
    
    // CSV数据行
    const csvRows = data.map(row => {
      return headers.map(header => {
        const value = row[header] || '';
        // 处理包含逗号、引号或换行符的值
        if (typeof value === 'string' && (value.includes(',') || value.includes('"') || value.includes('\n'))) {
          return `"${value.replace(/"/g, '""')}"`;
        }
        return value;
      }).join(',');
    });
    
    return [csvHeaders, ...csvRows].join('\n');
  };

  const downloadCSV = (csvContent: string, fileName: string) => {
    // 添加BOM以支持中文
    const BOM = '\uFEFF';
    const blob = new Blob([BOM + csvContent], { type: 'text/csv;charset=utf-8;' });
    
    // 创建下载链接
    const link = document.createElement('a');
    const url = URL.createObjectURL(blob);
    link.setAttribute('href', url);
    link.setAttribute('download', fileName);
    link.style.visibility = 'hidden';
    
    // 触发下载
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    
    // 清理URL对象
    URL.revokeObjectURL(url);
  };

  const onDownload = async (downloadParams: ListOrderReq) => {
    try {
      message.loading('正在生成文件...', 0);
      
      // 获取所有数据用于导出（不分页）
      const { data } = await listOrder({
        ...downloadParams,
        currentPage: 1,
        pageSize: 200000 // 获取大量数据
      })
      
      message.destroy();

      if (!data?.list || data.list.length === 0) {
        message.warning("没有数据可以下载")
        return
      }

      // 准备CSV数据
      const csvData = data.list.map((item, index) => ({
        '序号': index + 1,
        '系统订单号': item.orderId,
        '店铺订单号': item.partnerOrderId || '-',
        '商户订单号': item.merchantOrderId,
        '商户号': item.merchantId,
        '商户名称': item.merchantName,
        '订单金额': formatNumberWithCommasAndDecimals(item.amount),
        '实收金额': formatNumberWithCommasAndDecimals(item.receivedAmount),
        '订单状态': convertOrderStatus(item.status),
        '通知状态': convertNotifyStatus(item.notifyStatus),
        '下单时间': convertTimestamp(item.createAt),
        '支付类型': convertPayType(item.payType),
        '支付通道': convertChannelId(item.channelId),
        '下单账号': item.payAccount || '-',
        'skuId': item.skuId,
        '合作商号': item.partnerId,
        '合作商名称': item.partnerName,
        'IP地址': item.ip,
        '备注': item.remark || '-'
      }))

      // 定义CSV列头
      const headers = [
        '序号', '系统订单号', '店铺订单号', '商户订单号', '商户号', '商户名称',
        '订单金额', '实收金额', '订单状态', '通知状态', '下单时间', '支付类型',
        '支付通道', '下单账号', 'skuId', '合作商号', '合作商名称', 'IP地址', '备注'
      ];

      // 转换为CSV格式
      const csvContent = convertToCSV(csvData, headers);

      // 生成文件名
      const dateStr = downloadParams.startAt || new Date().toISOString().split('T')[0]
      const fileName = `交易记录_${dateStr}.csv`

      // 下载文件
      downloadCSV(csvContent, fileName);
      message.success("下载成功")
    } catch (e) {
      message.destroy();
      console.error(e);
      message.error("下载失败")
    }
  }

  return (
    <>
      <Card>
        <div style={{ display: 'Flex' }}>
          <TradingRecordSearchForm OnSearch={onSearch} OnDownload={onDownload} />
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
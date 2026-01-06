import { useCallback, useEffect, useState } from 'react';
import { Button, Card, Divider, message, Popconfirm, Space, Table } from 'antd';
import type { FormProps, TableProps } from 'antd';
import { useApis } from '../api/api';
import axios from 'axios';
import { GoodsCreateReq, IGoods, ListGoodsReq } from '../api/types';
import GoodsCreateModal from './modal/GoodsCreateModal';
import { convertTimestamp, GoodsStatus } from '../utils/biz';
import GoodsSearchForm from './searchform/GoodsSearchForm';
import { IGoodsSearchCondition } from './searchform/GoodsSearchForm';
import { PAGE_DEFAULT_INDEX, PAGE_SIZE } from '../components/types';
import { formatNumberWithCommasAndDecimals } from '../utils/utilb';

interface DataType extends IGoods {
  key: string;
}

const Goods = () => {
  const columns: TableProps<DataType>['columns'] = [
    {
      title: '合作商', dataIndex: 'partnerId', align: 'center', key: 'partnerId',
    },
    {
      title: '店铺', key: 'shopName', dataIndex: 'shopName', align: 'center',
    },
    {
      title: 'skuId', key: 'skuId', dataIndex: 'skuId', align: 'center',
    },
    {
      title: '商品金额', key: 'price', dataIndex: 'price', align: 'center', render: (_, d) => (
        formatNumberWithCommasAndDecimals(d.price)
      )
    },
    {
      title: '商品实际金额', key: 'realPrice', dataIndex: 'realPrice', align: 'center', render: (_, d) => (
        formatNumberWithCommasAndDecimals(d.realPrice)
      )
    },
    {
      title: '状态', key: 'status', dataIndex: 'status', align: 'center', render: (_, d) => { return (<span style={{ color: d.status === GoodsStatus.Enabled ? '#52c41a' : '#f5222d' }}> {d.status === GoodsStatus.Enabled ? '在售' : '非在售'}     </span>); }
    },
    {
      title: '创建时间', key: 'createAt', dataIndex: 'createAt', align: 'center', render: (_, d) => { return convertTimestamp(d.createAt) }
    },
    {
      title: '操作', key: 'action', fixed: 'right', // 固定最右边，配合Table的scroll={{ x: 'max-content' }}使用
      align: 'center', render: (_, d) => (
        <Space size="middle">
          <Button type="primary" size='small' onClick={() => openModal(d, true)}>修改</Button>
          <Popconfirm title="警告" description="请确认是否删除该商品？" onConfirm={() => deleteGoods(d.id)} >
            <Button type="primary" size='small' danger>删除</Button>
          </Popconfirm>
        </Space>),
    },
  ];

  const [list, setList] = useState<DataType[]>([])
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedData, setSelectedData] = useState<DataType>(null!);
  const [listLoadingStates, setListLoadingStates] = useState(false);
  let { listGoods, goodsDelete } = useApis()
  const [total, setTotal] = useState(0);
  const [reqParams, setReqParams] = useState<ListGoodsReq>({
    currentPage: PAGE_DEFAULT_INDEX,
    pageSize: PAGE_SIZE.TEN
  });

  const openModal = (selectedData: DataType | null = null, isOpen: boolean = false) => {
    setSelectedData(selectedData!)
    setIsModalOpen(isOpen);
  }

  const fetchListGoods = useCallback(async () => {
    try {
      setListLoadingStates(true)

      const { data } = await listGoods(reqParams)
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

  useEffect(() => {
    fetchListGoods()
  }, [reqParams])

  const onSuccess: FormProps<GoodsCreateReq>['onFinish'] = async () => {
    onSearch({});
    setIsModalOpen(false);
  }

  const deleteGoods = async (id: number) => {
    try {
      await goodsDelete({ id })
      onSearch({});
      message.success('删除成功');
    } catch (e) {
      if (axios.isAxiosError(e)) {
        let msg = e.response?.data?.message
        msg && message.success(msg);
      }
    }
  };

  const onSearch = (value: IGoodsSearchCondition) => {
    setReqParams({ ...value, currentPage: 1, pageSize: reqParams.pageSize })
  }

  return (
    <>
      <Card>
        <div style={{ display: 'flex' }}>
          <Button type="primary" onClick={() => { setIsModalOpen(true) }}>新增商品</Button>
          <Divider type="vertical" style={{ height: '32px', textAlign: 'center', alignContent: 'center', marginLeft: '20px', marginRight: '20px' }} />
          <GoodsSearchForm onSearch={onSearch} />
        </div>
        {/* <Divider /> */}
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
        {
          isModalOpen &&
          <GoodsCreateModal info={selectedData} isModalOpen={isModalOpen} onOk={onSuccess} onCancel={() => openModal()} />
        }
      </Card>
    </>
  )
}

export default Goods
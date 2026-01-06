import { useCallback, useEffect, useState } from 'react';
import { Table, Button, Space, Card, Divider, Popconfirm, message, Modal, QRCode, Typography } from 'antd';
import type { TableProps } from 'antd';
import { useApis } from '../api/api';
import { IMerchant, ListMerchantReq, MerchantUpdateReq, } from '../api/types';
import MerchantCreateModal from './modal/MerchantCreateModal';
import { convertEnable, convertTimestamp } from '../utils/biz';
import axios from 'axios';
import { formatNumberWithCommasAndDecimals, isEnable } from '../utils/utilb';
import { EnableStatus } from '../utils/constant';
import UpdateMerchantBalanceModal from './modal/UpdateMerchantBalanceModal';
import { useNavigate } from 'react-router-dom';
import MerchantSearchForm from './searchform/MerchantSearchForm';
import { MyPaginationConfig, PAGE_DEFAULT_INDEX, PAGE_SIZE } from '../components/types';

const { Text } = Typography;

interface DataType extends IMerchant {
  key: number;
}

function Merchant() {
  const columns: TableProps<DataType>['columns'] = [
    {
      title: '商户编号', dataIndex: 'id', key: 'id', align: 'center',
    },
    {
      title: '商户名称', dataIndex: 'nickname', key: 'nickname', align: 'center',
    },
    {
      title: '余额', key: 'balance', dataIndex: 'balance', align: 'center', render: (_, d) => {
        return <>
          <div>{formatNumberWithCommasAndDecimals(d.balance) || 0}</div>
          <span><Button type="link" size='small' onClick={() => { openBalanceModal(d.id, true) }}>调整</Button></span>
          <span><Button variant="link" color="danger" size='small' onClick={() => { goBalanceBill(d.id) }}>明细</Button></span>
        </>
      }
    },
    {
      title: '今日交易额', key: 'todayAmount', dataIndex: 'todayAmount', align: 'center', render: (_, d) => (
        formatNumberWithCommasAndDecimals(d.todayAmount)
      )
    },
    {
      title: '状态', key: 'enable', dataIndex: 'enable', align: 'center', render: (_, d) => {
        return (
          <span style={{ color: d.enable === EnableStatus.Enabled ? '#52c41a' : '#f5222d' }}>
            {convertEnable(d.enable)}
          </span>
        );
      }
    },
    {
      title: '创建时间', key: 'createAt', dataIndex: 'createAt', align: 'center', render: (_, d) => (
        convertTimestamp(d.createAt)
      )
    },
    {
      title: '秘钥', dataIndex: 'privateKey', key: 'privateKey', align: 'center',
    },
    {
      title: '归属', dataIndex: 'parentId', key: 'parentId', align: 'center',
    },
    {
      title: '备注', key: 'remark', dataIndex: 'remark', align: 'center', render: (_, d) => {
        return d.remark || '-'
      }
    },
    {
      title: '操作', key: 'action', align: 'center', fixed: 'right', width: 300, render: (_, d) => {
        let enable = isEnable(d.enable)
        return (
          <Space size='middle' wrap align="center" >
            <Button
              type="primary"
              size='small'
              danger={enable}
              onClick={
                () => enableMerchant(d.username, enable ? EnableStatus.Disabled : EnableStatus.Enabled)
              }>
              {enable ? '冻结' : '启用'}
            </Button>
            <Button type="primary" size='small' onClick={() => { openModal(d, true) }}>修改</Button>
            <Popconfirm title="警告" description="请确认是否删除该商户" onConfirm={() => deleteMerchant(d)} >
              <Button type="primary" size='small' danger >删除</Button>
            </Popconfirm>
            <Button type="primary" size='small' danger onClick={() => resetPassword(d.id)}>重置密码</Button>
            <Button type="primary" size='small' onClick={() => resetVerifiCode(d.id)}>重置验证码</Button>
            <Popconfirm
              title="验证码"
              icon={null}
              description={
                d?.urlKey ?
                  <QRCode value={d?.urlKey} size={320} /> :
                  <div>
                    <Text>无二维码，点击</Text>
                    <Text type="success">重置验证码</Text>
                  </div>
              }
              showCancel={false}
              okText="关闭"
            >
              <Button type="primary" size='small'>
                查看验证码
              </Button>
            </Popconfirm>
          </Space>
        )
      },
    },
  ];
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isBalanceModalOpen, setIsBalanceModalOpen] = useState(false);
  const [list, setList] = useState<DataType[]>([])
  let { merchantUpdate, merchantEnable, merchantResetPassword } = useApis()
  const [selectedData, setSelectedData] = useState<DataType>(null!);
  const [selectedId, setSelectedId] = useState<number>(null!);
  let navigate = useNavigate();

  const [listLoadingStates, setListLoadingStates] = useState(false);
  const [total, setTotal] = useState(0);
  let apis = useApis()

  const [reqParams, setReqParams] = useState<ListMerchantReq>({
    currentPage: PAGE_DEFAULT_INDEX,
    pageSize: PAGE_SIZE.TEN
  });

  const showModal = () => {
    setIsModalOpen(true);
  };

  const openModal = (selectedData: DataType | null = null, isOpen: boolean = false) => {
    setSelectedData(selectedData!)
    setIsModalOpen(isOpen);
  }

  const openBalanceModal = (selectedId: number | null = null, isOpen: boolean = false) => {
    setSelectedId(selectedId!)
    setIsBalanceModalOpen(isOpen);
  }

  const goBalanceBill = (id: number) => {
    navigate(`/admin/merchant/balanceBill?merchantId=${id}`);
  }

  const fetchListMerchant = useCallback(async () => {
    try {
      setListLoadingStates(true)

      const { data } = await apis.listMerchant(reqParams)
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
    fetchListMerchant()
  }, [reqParams])

  const onSuccess = async () => {
    fetchListMerchant();
    setIsModalOpen(false);
  }

  const onBalanceSuccess = async () => {
    fetchListMerchant();
    setIsBalanceModalOpen(false);
  }

  const deleteMerchant = async (data: DataType) => {
    try {
      let params: MerchantUpdateReq = {
        username: data.username,
        nickname: data.nickname,
        isDel: true
      }
      await merchantUpdate(params)
      fetchListMerchant()
      message.success('删除成功');
    } catch (e) {
      if (axios.isAxiosError(e)) {
        let msg = e.response?.data?.message
        msg && message.success(msg);
      }
    }
  };

  const enableMerchant = async (username: string, enable: number) => {
    try {
      await merchantEnable({ username, enable })
      fetchListMerchant()
      message.success(isEnable(enable) ? '启用成功' : '冻结成功')
    } catch (e) {
      if (axios.isAxiosError(e)) {
        let msg = e.response?.data?.message
        msg && message.success(msg);
      }
    }
  };

  const handleTableChange = (current: number, pageSize: number) => {
    setReqParams({ ...reqParams, currentPage: current, pageSize })
  };

  const onSearch = (value: ListMerchantReq) => {
    setReqParams({ ...value, currentPage: 1, pageSize: reqParams.pageSize })
  }

  const resetPassword = async (id: number) => {
    try {
      let { data } = await merchantResetPassword({ id })
      Modal.success({
        content: `重置密码成功, 密位为 ${data.password}`,
      });
    } catch (e) {
      if (axios.isAxiosError(e)) {
        let msg = e.response?.data?.message
        msg && message.error(msg);
      }
    }
  };

  const resetVerifiCode = async (id: number) => {
    try {
      await apis.merchantResetVerifiCode({ id })
      fetchListMerchant()
      message.success(`重置验证码成功, 查看二维码`)
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <>
      <Card>
        <div style={{ display: 'flex' }}>
          <Button type="primary" onClick={showModal}>新增</Button>
          <Divider type="vertical" style={{ height: '32px', textAlign: 'center', alignContent: 'center', marginLeft: '20px', marginRight: '20px' }} />
          <MerchantSearchForm onSearch={onSearch} />
        </div>
        <Table<DataType>
          bordered
          size='middle'
          columns={columns}
          dataSource={list}
          scroll={{ x: 'max-content' }}
          loading={listLoadingStates}
          pagination={{
            current: reqParams.currentPage,
            pageSize: reqParams.pageSize,
            total: total,
            onChange: handleTableChange,
          } as MyPaginationConfig} />

        {
          isModalOpen &&
          <MerchantCreateModal info={selectedData} isModalOpen={isModalOpen} onOk={onSuccess} onCancel={() => openModal()} />
        }

        {
          isBalanceModalOpen &&
          <UpdateMerchantBalanceModal merchantId={selectedId} isModalOpen={isBalanceModalOpen} onOk={onBalanceSuccess} onCancel={() => openBalanceModal()} />
        }
      </Card>
    </>
  )
}

export default Merchant
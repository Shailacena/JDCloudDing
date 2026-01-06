import { useState, useEffect, useCallback } from 'react';
import { Space, Table, Button, message, Card, Divider, Popconfirm, Modal, QRCode, Typography } from 'antd';
import type { TableProps } from 'antd';
import { useApis } from '../api/api';
import { IPartner, ListPartnerReq } from '../api/types';
import { useAppContext } from '../AppProvider';
import PartnerSearchForm from './searchform/PartnerSearchForm';
import PartnerCreateModal from './modal/PartnerCreateModal';
import axios from 'axios';
import { EnableStatus } from '../utils/constant';
import UpdatePartnerBalanceModal from './modal/UpdatePartnerBalanceModal';
import { PAGE_DEFAULT_INDEX, PAGE_SIZE } from '../components/types';
import { useNavigate } from 'react-router-dom';
import { AllPartnerType, convertEnable, isJDShop, PartnerType } from '../utils/biz';
import { formatNumberWithCommasAndDecimals, isEnable, toPercentWithFixed } from '../utils/utilb';

const { Text } = Typography;

interface DataType extends IPartner {
  key: number;
}

enum ActionType {
  ENABLE,
  UPDATE,
  SYN_DATA,
  TO_GOODS,
  DELETE,
  CREDIT,
  RESETPASSWORD
};

function Partner() {
  const columns: TableProps<DataType>['columns'] = [
    {
      title: 'ID', dataIndex: 'id', key: 'id', align: 'center',
    },
    {
      title: '合作商名称', dataIndex: 'nickname', key: 'nickname', align: 'center',
    },
    {
      title: '合作商类型', dataIndex: 'type', key: 'type', align: 'center', render: (_, d) => {
        return AllPartnerType.find((item) => d.type === item.id)?.label
      }
    },
    {
      title: '押金', key: 'balance', dataIndex: 'balance', align: 'center', render: (_, d) => {
        return <div className='text-center'>
          <div>{formatNumberWithCommasAndDecimals(d.balance) || 0}</div>
          <div><Button type="link" size='small' onClick={() => { openBalanceModal(d.id, true) }}>调整</Button></div>
        </div>
      }
    },
    {
      title: '今日订单总额', key: 'todayOrderAmount', dataIndex: 'todayOrderAmount', align: 'center', render: (_, d) => (
        formatNumberWithCommasAndDecimals(d.todayOrderAmount)
      )
    },
    {
      title: '今日成功订单总额', key: 'todaySuccessAmount', dataIndex: 'todaySuccessAmount', align: 'center', render: (_, d) => (
        formatNumberWithCommasAndDecimals(d.todaySuccessAmount)
      )
    },
    {
      title: '今日订单总数', key: 'todayOrderNum', dataIndex: 'todayOrderNum', align: 'center',
    },
    {
      title: '今日成功订单总数', key: 'todaySuccessOrderNum', dataIndex: 'todaySuccessOrderNum', align: 'center',
    },
    {
      title: '今日成功率', key: 'todaySuccessOrderRate', dataIndex: 'todaySuccessOrderRate', align: 'center', render: (_, d) => (
        toPercentWithFixed(d.todaySuccessOrderNum, d.todayOrderNum) + "%"
      )
    },
    {
      title: '近1小时成功率', align: 'center', render: (_, d) => (
        toPercentWithFixed(d.last1HourSuccess, d.last1HourTotal) + "%"
      )
    },
    {
      title: '近30分钟成功率', align: 'center', render: (_, d) => (
        toPercentWithFixed(d.last30MinutesSuccess, d.last30MinutesTotal) + "%"
      )
    },
    {
      title: '状态', key: 'enable', dataIndex: 'enable', align: 'center', render: (_, d) => (
        <span style={{ color: d.enable === EnableStatus.Enabled ? '#52c41a' : '#f5222d' }}>
          {convertEnable(d.enable)}
        </span>
      )
    },
    {
      title: '优先级', key: 'priority', dataIndex: 'priority', align: 'center',
    },
    {
      title: '归属', dataIndex: 'parentId', key: 'parentId', align: 'center',
    },
    // {
    //   title: '上级代理', key: 'superiorAgent', dataIndex: 'superiorAgent', align: 'center', render: (_, d) => {
    //     return d.superiorAgent || '-'
    //   }
    // },
    // {
    //   title: '等级', key: 'level', dataIndex: 'level', align: 'center',
    // },
    {
      title: '操作', key: 'action', fixed: 'right', align: 'left', width: 300, render: (_, d) => (
        <Space direction="vertical">
          <Space size='middle' direction="horizontal">
            <Button type="primary" size='small' onClick={() => handleUpdate(ActionType.ENABLE, d)} danger={isEnable(d.enable)}>{isEnable(d.enable) ? '冻结' : '启用'}</Button>
            <Button type="primary" size='small' onClick={() => handleUpdate(ActionType.UPDATE, d)}>修改</Button>
            <Popconfirm title="警告" description="请确认是否删除该合作商和该合作商的商品？" onConfirm={() => handleUpdate(ActionType.DELETE, d)} >
              <Button type="primary" size='small' danger>删除</Button>
            </Popconfirm>
          </Space>
          <Space size='middle' direction="horizontal">
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
            <Button type="primary" size='small' danger onClick={() => handleUpdate(ActionType.RESETPASSWORD, d)}>重置密码</Button>
            <Button type="primary" size='small' onClick={() => resetVerifiCode(d.id)}>重置验证码</Button>
          </Space>
          <Space size='middle' direction="horizontal">
            <Button type="primary" size='small' onClick={() => handleUpdate(ActionType.TO_GOODS, d)}>查看商品</Button>
            {
              !isJDShop(d.channelId) && <Button type="primary" size='small' onClick={() => handleUpdate(ActionType.SYN_DATA, d)}
                loading={
                  //@ts-ignore
                  loadingStates[d.id]}>同步商品</Button>
            }
            {
              d.type === PartnerType.Anssy && <Button type="primary" size='small' onClick={() => anssyAuth(d)}>店铺授权</Button>
            }
          </Space>
        </Space>
      ),
    },
  ];

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isBalanceModalOpen, setIsBalanceModalOpen] = useState(false);
  const [list, setList] = useState<DataType[]>([])
  const [loadingStates, setLoadingStates] = useState({});
  const [listLoadingStates, setListLoadingStates] = useState(false);
  const [selectedData, setSelectedData] = useState<DataType>(null!);
  const [selectedId, setSelectedId] = useState<number>(null!);

  const [total, setTotal] = useState(0);
  const [reqParams, setReqParams] = useState<ListPartnerReq>({
    currentPage: PAGE_DEFAULT_INDEX,
    pageSize: PAGE_SIZE.TEN
  });

  let navigate = useNavigate();
  let ctx = useAppContext();
  let apis = useApis()

  const openBalanceModal = (selectedId: number | null = null, isOpen: boolean = false) => {
    setSelectedId(selectedId!)
    setIsBalanceModalOpen(isOpen);
  }

  const anssyAuth = (value: DataType) => {
    window.open(`https://tao.anssy.com/msgpush/taobao/auth?state=${value.id},7`, '_blank');
  }

  const fetchListPartner = useCallback(async () => {
    try {
      setListLoadingStates(true)

      const { data } = await apis.listPartner(reqParams)
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
    fetchListPartner()
  }, [reqParams])


  const onSuccess = () => {
    fetchListPartner()
    setIsModalOpen(false)
  };


  const handleUpdate = async (type: ActionType, value: DataType) => {
    try {
      switch (type) {
        case ActionType.ENABLE: {
          await apis.partnerUpdate({
            id: value.id,
            type: value.type,
            nickname: value.nickname,
            priority: value.priority,
            aqsAppSecret: value.aqsAppSecret,
            aqsToken: value.aqsToken,
            enable: isEnable(value.enable) ? EnableStatus.Disabled : EnableStatus.Enabled,
            darkNumberLength: value.darkNumberLength
          })
          fetchListPartner()
          if (!isEnable(value.enable)) {
            message.success('启用成功')
          } else if (isEnable(value.enable)) {
            message.success('冻结成功')
          }
          break;
        }
        case ActionType.RESETPASSWORD: {
          let { data } = await apis.partnerResetPassword(value)
          Modal.success({
            content: `重置密码成功, 密位为 ${data.password}`,
          });
          break;
        }
        case ActionType.DELETE: {
          await apis.partnerDelete(value)
          fetchListPartner()
          ctx.partnerList = [];
          message.success('删除成功');
          break;
        }
        case ActionType.UPDATE: {
          openModal(value, true)
          break;
        }
        case ActionType.SYN_DATA: {
          setLoadingStates(() => ({ [value.id]: true }));
          try {
            await apis.partnerSyncGoods({ id: value.id })
            message.success('同步成功');
          } finally {
            setLoadingStates(() => ({ [value.id]: false }));
          }
          break;
        }
        case ActionType.TO_GOODS: {
          navigate(`/admin/partner/goods?partnerId=${value.id}`);
          break;
        }
      }
    } catch (e) {
      if (axios.isAxiosError(e)) {
        let msg = e.response?.data?.message
        msg && message.error(msg)
      }
    }
  }

  const handleTableChange = (current: number, pageSize: number) => {
    setReqParams({ ...reqParams, currentPage: current, pageSize })
  };

  const onSearch = (value: ListPartnerReq) => {
    setReqParams({ ...value, currentPage: 1, pageSize: reqParams.pageSize })
  }

  const openModal = (selectedData: DataType | null = null, isOpen: boolean = false) => {
    setSelectedData(selectedData!)
    setIsModalOpen(isOpen)
  }

  const onBalanceSuccess = async () => {
    fetchListPartner();
    setIsBalanceModalOpen(false);
  }

  const resetVerifiCode = async (id: number) => {
    try {
      await apis.partnerResetVerifiCode({ id })
      fetchListPartner()
      message.success(`重置验证码成功, 查看二维码`)
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <>
      <Card>
        <div style={{ display: 'flex' }}>
          <Button type="primary" onClick={() => { setIsModalOpen(true) }} >新增</Button>
          <Divider type="vertical" style={{ height: '32px', textAlign: 'center', alignContent: 'center', marginLeft: '20px', marginRight: '20px' }} />
          <PartnerSearchForm onSearch={onSearch} />
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
        {
          isModalOpen && <PartnerCreateModal info={selectedData} isModalOpen={isModalOpen} onOk={onSuccess} onCancel={() => openModal()} />
        }

        {
          isBalanceModalOpen &&
          <UpdatePartnerBalanceModal partnerId={selectedId} isModalOpen={isBalanceModalOpen} onOk={onBalanceSuccess} onCancel={() => openBalanceModal()} />
        }
      </Card>
    </>
  )
}

export default Partner
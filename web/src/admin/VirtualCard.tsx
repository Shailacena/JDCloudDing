import { useCallback, useEffect, useState } from 'react';
import { Button, Card, DatePicker, Form, Input, message, Modal, Space, Table, Tag } from 'antd';
import type { FormProps, TableProps } from 'antd';
import { SearchOutlined, CloseCircleOutlined } from '@ant-design/icons';
import { useApis } from '../api/api';
import { IPriceCard, ListCardReq, VirtualCardGenerateReq } from '../api/types';
import { PAGE_DEFAULT_INDEX, PAGE_SIZE } from '../components/types';
import { formatNumberWithCommasAndDecimals } from '../utils/utilb';
import dayjs from 'dayjs';

const { RangePicker } = DatePicker;

interface DataType extends IPriceCard {
  key: string;
}

const VirtualCardPage = () => {
  const [list, setList] = useState<DataType[]>([])
  const [listLoadingStates, setListLoadingStates] = useState(false);
  const [total, setTotal] = useState(0);
  const [reqParams, setReqParams] = useState<ListCardReq>({
    currentPage: PAGE_DEFAULT_INDEX,
    pageSize: PAGE_SIZE.TEN,
  });

  const [isGenerateModalOpen, setIsGenerateModalOpen] = useState(false);
  const [generateForm] = Form.useForm();

  const apis = useApis()

  const fetchList = useCallback(async () => {
    setListLoadingStates(true)
    try {
      const res = await apis.listVirtualCard(reqParams)
      if (res?.data) {
        setList(res.data.list?.map((item: IPriceCard) => ({ ...item, key: String(item.id) })) || [])
        setTotal(res.data.total || 0)
      }
    } catch (e) {
      console.error(e);
    } finally {
      setListLoadingStates(false)
    }
  }, [reqParams, apis])

  useEffect(() => {
    fetchList()
  }, [reqParams])

  const handleTableChange = (current: number, pageSize: number) => {
    setReqParams({ ...reqParams, currentPage: current, pageSize })
  };

  const handleDelete = async (ids: number[]) => {
    try {
      await apis.cardDelete({ ids })
      message.success('删除成功')
      fetchList()
    } catch (e) {
      console.error(e);
    }
  };

  const handleDeleteByCondition = async () => {
    try {
      await Modal.confirm({
        title: '确认删除',
        content: '确定要删除所有符合搜索条件的卡密吗？此操作不可恢复！',
        okButtonProps: { danger: true },
        onOk: async () => {
          const res = await apis.cardDeleteByCondition({
            ...reqParams,
            cardType: 'virtual',
          })
          message.success(`删除成功，共删除 ${res?.data?.count || 0} 条`)
          fetchList()
        },
      })
    } catch (e) {
      console.error(e);
    }
  };

  const columns: TableProps<DataType>['columns'] = [
    {
      title: 'ID', dataIndex: 'id', key: 'id', align: 'center', width: 60,
    },
    {
      title: '卡号', dataIndex: 'cardNo', key: 'cardNo', align: 'center',
    },
    {
      title: '密码', dataIndex: 'password', key: 'password', align: 'center',
    },
    {
      title: '卡组', dataIndex: 'cardGroup', key: 'cardGroup', align: 'center',
    },
    {
      title: '面额', dataIndex: 'amount', key: 'amount', align: 'center', render: (_, d) => (
        formatNumberWithCommasAndDecimals(d.amount)
      )
    },
    {
      title: '批次', dataIndex: 'batchNo', key: 'batchNo', align: 'center',
    },
    {
      title: '状态', dataIndex: 'usedStatus', key: 'usedStatus', align: 'center', render: (_, d) => (
        <Tag color={d.usedStatus ? 'green' : 'blue'}>
          {d.usedStatus ? '已使用' : '未使用'}
        </Tag>
      )
    },
    {
      title: '订单ID', dataIndex: 'orderId', key: 'orderId', align: 'center', render: (_, d) => d.orderId || '-',
    },
    {
      title: '创建时间', dataIndex: 'createAt', key: 'createAt', align: 'center', render: (_, d) => dayjs(d.createAt * 1000).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '操作', key: 'action', fixed: 'right', align: 'center', render: (_, d) => (
        <Button type="primary" size='small' danger onClick={() => handleDelete([d.id])}>删除</Button>
      )
    },
  ];

  const handleSearch = (values: any) => {
    setReqParams({
      ...reqParams,
      ...values,
      startTime: values.dateRange?.[0]?.format('YYYY-MM-DD'),
      endTime: values.dateRange?.[1]?.format('YYYY-MM-DD'),
      currentPage: 1,
    })
  };

  const handleReset = () => {
    setReqParams({
      currentPage: PAGE_DEFAULT_INDEX,
      pageSize: PAGE_SIZE.TEN,
    })
  };

  const handleGenerate: FormProps<VirtualCardGenerateReq>['onFinish'] = async (values) => {
    try {
      await apis.cardGenerateVirtual({
        ...values,
        cardNoLen: Number(values.cardNoLen),
        passwordLen: Number(values.passwordLen),
        amount: Number(values.amount),
        count: Number(values.count),
      })
      message.success('生成成功');
      setIsGenerateModalOpen(false);
      generateForm.resetFields();
      fetchList();
    } catch (e) {
      console.error(e);
      message.error('生成失败');
    }
  };

  return (
    <>
      <Card>
        <Form
          layout="inline"
          onFinish={handleSearch}
          style={{ marginBottom: 16, flexWrap: 'wrap', gap: '8px 16px' }}
        >
          <Form.Item name="cardNo" label="卡号">
            <Input placeholder="请输入卡号" style={{ width: 150 }} />
          </Form.Item>
          <Form.Item name="cardGroup" label="卡组">
            <Input placeholder="请输入卡组" style={{ width: 150 }} />
          </Form.Item>
          <Form.Item name="batchNo" label="批次">
            <Input placeholder="请输入批次" style={{ width: 150 }} />
          </Form.Item>
          <Form.Item name="dateRange" label="时间">
            <RangePicker />
          </Form.Item>
          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit" icon={<SearchOutlined />}></Button>
              <Button type="primary" onClick={handleReset} icon={<CloseCircleOutlined />}></Button>
              <Button type="primary" danger onClick={handleDeleteByCondition}>删除所搜</Button>
            </Space>
          </Form.Item>
        </Form>

        <div style={{ marginBottom: 16 }}>
          <Button type="primary" onClick={() => setIsGenerateModalOpen(true)}>生成卡密</Button>
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
          }}
        />
      </Card>

      <Modal title="生成虚拟卡密" open={isGenerateModalOpen} onCancel={() => setIsGenerateModalOpen(false)} footer={null}>
        <Form
          form={generateForm}
          layout="vertical"
          onFinish={handleGenerate}
        >
          <Form.Item name="prefix" label="卡号前缀" rules={[{ required: true, message: '请输入卡号前缀' }]}>
            <Input placeholder="例如: CARD" />
          </Form.Item>
          <Form.Item name="cardNoLen" label="卡号长度" rules={[{ required: true, message: '请输入卡号长度' }, { validator: (_, value) => {
            if (!value || isNaN(Number(value))) {
              return Promise.reject('请输入有效的数字');
            }
            if (Number(value) < 12) {
              return Promise.reject('长度不能低于12位');
            }
            return Promise.resolve();
          } }]}>
            <Input type="number" placeholder="最低12位" />
          </Form.Item>
          <Form.Item name="passwordLen" label="密码长度" rules={[{ required: true, message: '请输入密码长度' }, { validator: (_, value) => {
            if (!value || isNaN(Number(value))) {
              return Promise.reject('请输入有效的数字');
            }
            if (Number(value) < 1) {
              return Promise.reject('密码长度不能少于1位');
            }
            return Promise.resolve();
          } }]}>
            <Input type="number" placeholder="密码长度" />
          </Form.Item>
          <Form.Item name="cardGroup" label="卡组" rules={[{ required: true, message: '请输入卡组' }]}>
            <Input placeholder="例如: 联通100" />
          </Form.Item>
          <Form.Item name="amount" label="面额" rules={[{ required: true, message: '请输入面额' }, { validator: (_, value) => {
            if (!value || isNaN(Number(value))) {
              return Promise.reject('请输入有效的数字');
            }
            return Promise.resolve();
          } }]}>
            <Input type="number" placeholder="例如: 100.00" />
          </Form.Item>
          <Form.Item name="count" label="生成数量" rules={[{ required: true, message: '请输入生成数量' }, { validator: (_, value) => {
            if (!value || isNaN(Number(value))) {
              return Promise.reject('请输入有效的数字');
            }
            if (Number(value) < 1) {
              return Promise.reject('数量至少为1');
            }
            return Promise.resolve();
          } }]}>
            <Input type="number" placeholder="例如: 100" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              生成
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}

export default VirtualCardPage

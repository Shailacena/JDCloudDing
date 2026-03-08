import { useCallback, useEffect, useState } from 'react';
import { Button, Card, DatePicker, Form, Input, message, Modal, Space, Table, Tag, Upload } from 'antd';
import type { FormProps, TableProps } from 'antd';
import { useApis } from '../api/api';
import * as XLSX from 'xlsx';
import { IPriceCard, ListCardReq, CardInfo } from '../api/types';
import { PAGE_DEFAULT_INDEX, PAGE_SIZE } from '../components/types';
import { formatNumberWithCommasAndDecimals } from '../utils/utilb';
import { UploadOutlined, SearchOutlined, CloseCircleOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';

const { RangePicker } = DatePicker;
const { TextArea } = Input;

interface DataType extends IPriceCard {
  key: string;
}

const RealCardPage = () => {
  const [list, setList] = useState<DataType[]>([])
  const [listLoadingStates, setListLoadingStates] = useState(false);
  const [total, setTotal] = useState(0);
  const [reqParams, setReqParams] = useState<ListCardReq>({
    currentPage: PAGE_DEFAULT_INDEX,
    pageSize: PAGE_SIZE.TEN,
  });

  const [isImportModalOpen, setIsImportModalOpen] = useState(false);
  const [importCards, setImportCards] = useState<CardInfo[]>([]);

  const apis = useApis()

  const fetchList = useCallback(async () => {
    setListLoadingStates(true)
    try {
      const res = await apis.listRealCard(reqParams)
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
            cardType: 'real',
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

  const beforeUpload = (file: File) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      try {
        const workbook = XLSX.read(e.target?.result as ArrayBuffer, { type: 'array' });
        const firstSheet = workbook.Sheets[workbook.SheetNames[0]];
        const jsonData = XLSX.utils.sheet_to_json<CardInfo>(firstSheet);

        const cards: CardInfo[] = jsonData.map(item => ({
          cardNo: String(item.cardNo) || '',
          password: String(item.password) || '',
          cardGroup: String(item.cardGroup) || '',
          amount: Number(item.amount) || 0,
        })).filter(item => item.cardNo && item.password);

        setImportCards(cards)
        message.success(`解析成功，共 ${cards.length} 条数据`);
      } catch (error) {
        message.error('文件解析失败');
      }
    };
    reader.readAsArrayBuffer(file);
    return false;
  };

  const handleImport = async () => {
    if (importCards.length === 0) {
      message.warning('请先上传Excel文件');
      return;
    }
    try {
      await apis.cardCreate({ cards: importCards })
      message.success('导入成功');
      setIsImportModalOpen(false);
      setImportCards([]);
      fetchList();
    } catch (e) {
      console.error(e);
      message.error('导入失败');
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
          <Space>
            <Button type="primary" onClick={() => setIsImportModalOpen(true)}>导入卡密</Button>
            <Button onClick={() => {
              const sampleData = [
                { cardNo: 'CARD123456789', password: 'PASSWORD1', cardGroup: 'VIP卡组', amount: 100 },
                { cardNo: 'CARD987654321', password: 'PASSWORD2', cardGroup: 'VIP卡组', amount: 200 },
              ];
              const worksheet = XLSX.utils.json_to_sheet(sampleData);
              const workbook = XLSX.utils.book_new();
              XLSX.utils.book_append_sheet(workbook, worksheet, '卡密导入模板');
              XLSX.writeFile(workbook, '卡密导入模板.xlsx');
              message.success('示例文件已导出');
            }}>导出示例Excel</Button>
          </Space>
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

      <Modal title="导入真实卡密" open={isImportModalOpen} onCancel={() => setIsImportModalOpen(false)} footer={null}>
        <Upload
          accept=".xlsx,.xls"
          beforeUpload={beforeUpload}
          showUploadList={false}
        >
          <Button icon={<UploadOutlined />}>选择Excel文件</Button>
        </Upload>
        <div style={{ marginTop: 16, marginBottom: 16 }}>
          <TextArea rows={4} value={importCards.map(c => `${c.cardNo},${c.password},${c.cardGroup},${c.amount}`).join('\n')} placeholder="Excel解析结果预览" />
        </div>
        <div style={{ textAlign: 'center' }}>
          <Button type="primary" onClick={handleImport} disabled={importCards.length === 0}>
            导入 ({importCards.length} 条)
          </Button>
        </div>
      </Modal>
    </>
  )
}

export default RealCardPage

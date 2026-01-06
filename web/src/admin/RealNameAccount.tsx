import { useCallback, useEffect, useState } from 'react';
import { Modal, Form, Table, Input, Button, Card, Divider, message, Upload } from 'antd';
import type { FormProps, TableProps } from 'antd';
import { useApis } from '../api/api';
import { BaseRealNameAccount, IRealNameAccount, ListRealNameAccountReq, RealNameAccountCreateReq } from '../api/types';
import axios from 'axios';
import { EnableStatus } from '../utils/constant';
import { convertEnable } from '../utils/biz';
import { PAGE_DEFAULT_INDEX, PAGE_SIZE } from '../components/types';
import * as XLSX from 'xlsx';
import { UploadOutlined } from '@ant-design/icons';

const { TextArea } = Input;

interface DataType extends IRealNameAccount {
  key: number
}

type FieldType = {
  accounts: string
};

interface ExcelData {
  name: string;
  idCard: string;
  phone: string;
  address: string;
  [key: string]: any;
}

const columns: TableProps<DataType>['columns'] = [
  {
    title: 'ID',
    dataIndex: 'idNumber',
    key: 'idNumber',
  },
  {
    title: '名称',
    dataIndex: 'name',
    key: 'name',
  },
  {
    title: '手机',
    dataIndex: 'mobile',
    key: 'mobile',
  },
  {
    title: '地址',
    dataIndex: 'address',
    key: 'address',
  },
  {
    title: '状态',
    key: 'enable',
    dataIndex: 'enable',
    render: (_, d) => (
      <span style={{ color: d.enable === EnableStatus.Enabled ? '#52c41a' : '#f5222d' }}>
        {convertEnable(d.enable)}
      </span>
    )
  },
];

function RealNameAccount() {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [list, setList] = useState<DataType[]>([])
  let { listRealNameAccount, realNameAccountCreate } = useApis()
  const [listLoadingStates, setListLoadingStates] = useState(false);
  const [total, setTotal] = useState(0);
  const [reqParams, setReqParams] = useState<ListRealNameAccountReq>({
    currentPage: PAGE_DEFAULT_INDEX,
    pageSize: PAGE_SIZE.TEN
  });

  const [form] = Form.useForm();

  const showModal = () => {
    setIsModalOpen(true);
  };

  const handleOk = () => {
    setIsModalOpen(false);
  };

  const handleCancel = () => {
    setIsModalOpen(false);
  };

  const onFinish: FormProps<FieldType>['onFinish'] = async (value) => {
    try {
      let accountList: BaseRealNameAccount[] = []
      if (value.accounts) {
        let list = value.accounts.split(/[(\r\n)\r\n]+/)

        list?.forEach((line: string) => {
          line = line.trim().replace(/\s+/, ",")
          let accounts = line.split(",")
          if (accounts.length > 1) {
            accountList.push({
              idNumber: accounts[0],
              name: accounts[1],
              mobile: accounts[2],
              address: accounts[3]
            })
          }
        })
      }

      let data: RealNameAccountCreateReq = {
        accountList: accountList,
        remark: ""
      }

      await realNameAccountCreate(data)

      onSearch(reqParams)
      setIsModalOpen(false);
      message.success('导入成功')
    } catch (e) {
      if (axios.isAxiosError(e)) {
        let msg = e.response?.data?.message
        msg && message.error(msg)
      }
    }
  };

  const onSearch = (value: ListRealNameAccountReq) => {
    setReqParams({ ...value, currentPage: 1, pageSize: reqParams.pageSize })
  }

  const fetchListRealNameAccount = useCallback(async () => {
    try {
      setListLoadingStates(true)

      const { data } = await listRealNameAccount(reqParams)
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
    fetchListRealNameAccount()
  }, [reqParams])

  const handleTableChange = (current: number, pageSize: number) => {
    setReqParams({ ...reqParams, currentPage: current, pageSize })
  };

  const beforeUpload = (file: File) => {
    const reader = new FileReader();
    reader.onload = (e) => {
      try {
        const workbook = XLSX.read(e.target?.result as ArrayBuffer, { type: 'array' });
        const firstSheet = workbook.Sheets[workbook.SheetNames[0]];
        const jsonData = XLSX.utils.sheet_to_json<ExcelData>(firstSheet);

        const formatData = (arr: Array<Record<string, string>>): string => {
          return arr.map(item =>
            `${item.身份证},${item.姓名},${item.电话},${item.地址}`
          ).join('\n');
        };

        console.log(formatData(jsonData))

        // setTextAreaValue(formatData(jsonData))
        // setTextAreaValue('123')
        form.setFieldsValue({ accounts: formatData(jsonData) })

        message.success('文件解析成功');
      } catch (error) {
        message.error('文件解析失败');
      }
    };
    reader.readAsArrayBuffer(file);
    return false; // 阻止自动上传
  };

  return (
    <>
      <Card>
        <div>
          <Button type="primary" onClick={showModal}>批量导入实名资料</Button>
        </div>
        <Divider />
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

        <Modal title="导入实名资料" destroyOnClose open={isModalOpen} onOk={handleOk} onCancel={handleCancel} footer={null}>

          <Upload
            accept=".xlsx,.xls"
            beforeUpload={beforeUpload}
            showUploadList={false}
          >
            <Button icon={<UploadOutlined />}>选择Excel文件</Button>
          </Upload>

          <Divider />
          <Form
            form={form}
            name="basic"
            autoComplete="off"
            onFinish={onFinish}
          >
            <Form.Item<FieldType>
              name="accounts"
              label="账号"
            >
              <TextArea rows={4} />
            </Form.Item>

            <Form.Item>
              <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                <Button size="large" block type="primary" htmlType="submit" style={{ width: 100 }}>
                  提交
                </Button>
              </div>
            </Form.Item>
          </Form >
        </Modal>
      </Card>
    </>
  )
}

export default RealNameAccount
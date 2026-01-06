import { SearchOutlined } from '@ant-design/icons';
import { Table, Button, Form, message, DatePicker, Card } from 'antd';
import type { SelectProps, TableProps } from 'antd';
import CurrentLocation from '../components/CurrentLocation';
import { getRouteConfig } from './RouteConfigs';
import { getDataFormat } from '../utils/Tool';
import { useEffect, useState } from 'react';

interface DataType {
  key: string;
}

const columns: TableProps<DataType>['columns'] = [
  {
    title: '备注', dataIndex: 'notes', key: 'notes', align: 'center',
  },
  {
    title: '变更金额', key: 'changeMoney', dataIndex: 'changeMoney', align: 'center',
  },
  {
    title: '当前余额', key: 'money', dataIndex: 'money', align: 'center',
  },
  {
    title: '时间', key: 'createAt', dataIndex: 'createAt', align: 'center', render: (text) => {
      const date = new Date(text);
      return getDataFormat(date);
    }
  }
];

function CashFlow() {
  const [list, _] = useState<DataType[]>([])
  // let { listPartner1Bill } = useApis()

  const fetchListPartnerBill = async () => {
    // const { data } = await listPartner1Bill()
    // let d: DataType[] = data?.list?.map((item, index) => {
    //   let newItem: DataType = {
    //     key: index.toString(),
    //     ...item
    //   }
    //   return newItem
    // })
    // setList(d)
  }

  useEffect(() => {
    fetchListPartnerBill()
  }, [])

  return (
    <>
      <div style={{ marginBottom: '10px' }}>
        <CurrentLocation routeconfigs={getRouteConfig()} />
      </div>
      <Card>
        <div style={{ display: 'Flex' }}>
          <SearchForm />
          <Button type="primary" onClick={() => toggleModal()}>
            导出
          </Button>
        </div>
        <Table<DataType>
          bordered
          size='small'
          pagination={{ pageSize: 12 }}
          scroll={{ x: 'max-content' }}
          columns={columns}
          dataSource={list || []} />
      </Card>
    </>
  )
}

const options: SelectProps['options'] = [];

for (let i = 10; i < 36; i++) {
  options.push({
    value: i.toString(36) + i,
    label: i.toString(36) + i,
  });
}

const SearchForm = () => {
  const [form] = Form.useForm();

  const onFinish = (_: any) => {
    message.success('Search Success!');
  };

  const { RangePicker } = DatePicker;

  return (
    <Form
      form={form}
      layout="inline"
      onFinish={onFinish}
      style={{ marginBottom: 16 }}
    >
      <Form.Item
        name="searchKeyword"
        label="Date"
      >
        <RangePicker
          id={{
            start: 'startInput',
            end: 'endInput',
          }} />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit" icon={<SearchOutlined />}>
        </Button>
      </Form.Item>
    </Form>
  );
};

const toggleModal = () => {
  message.warning('功能还未完成...');
};

export default CashFlow
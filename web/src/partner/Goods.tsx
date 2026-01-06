import { Table, Button, Form, message, Input, Card } from 'antd';
import type { TableProps } from 'antd';
import { SearchOutlined } from '@ant-design/icons';
import CurrentLocation from '../components/CurrentLocation';
import { getRouteConfig } from './RouteConfigs';
import { useCallback, useEffect, useState } from 'react';
import { useApis } from '../api/api';
import { useAppContext } from '../AppProvider';
import { IGoods, ListGoodsReq } from '../api/types';
import { convertTimestamp, GoodsStatus } from '../utils/biz';
import { PAGE_DEFAULT_INDEX, PAGE_SIZE } from '../components/types';
import { formatNumberWithCommasAndDecimals } from '../utils/utilb';

interface DataType extends IGoods {
  key: string;
}

interface IGoodsSearchCondition {
  skuId?: string;
}

interface SearchBarProps {
  onSearch: (conditions: IGoodsSearchCondition) => void;
}

const SearchForm: React.FC<SearchBarProps> = ({ onSearch }) => {
  const [form] = Form.useForm<IGoodsSearchCondition>();

  const onFinish = (values: IGoodsSearchCondition = {}) => {
    onSearch(values);
  };

  return (
    <Form
      form={form}
      layout="inline"
      onFinish={onFinish}
      style={{ marginBottom: 16 }}
    >
      <Form.Item<IGoodsSearchCondition> name="skuId">
        <Input placeholder="商品sku" allowClear />
      </Form.Item>

      <Form.Item>
        <Button type="primary" htmlType="submit" icon={<SearchOutlined />}></Button>
      </Form.Item>
    </Form>
  );
};

function Goods() {
  const columns: TableProps<DataType>['columns'] = [
    {
      title: '合作商', dataIndex: 'partnerId', align: 'center', key: 'partnerId',
    },
    {
      title: 'sku', key: 'skuId', dataIndex: 'skuId', align: 'center',
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
      title: '状态', key: 'status', dataIndex: 'status', align: 'center', render: (_, d) => {
        return (
          <span style={{ color: d.status === GoodsStatus.Enabled ? '#52c41a' : '#f5222d' }}>
            {d.status === GoodsStatus.Enabled ? '在售' : '非在售'}
          </span>
        );
      }
    },
    {
      title: '创建时间', key: 'createAt', dataIndex: 'createAt', align: 'center', render: (_, d) => {
        return convertTimestamp(d.createAt)
      }
    },
  ];

  const [list, setList] = useState<DataType[]>([])
  let apis = useApis()
  let ctx = useAppContext();
  const [total, setTotal] = useState(0);
  const [listLoadingStates, setListLoadingStates] = useState(false);
  const [syncLoadingStates, setsyncLoadingStates] = useState(false);
  const [reqParams, setReqParams] = useState<ListGoodsReq>({
    currentPage: PAGE_DEFAULT_INDEX,
    pageSize: PAGE_SIZE.TEN,
    partnerId: ctx.cookie.id
  });

  const fetchListGoods = useCallback(async () => {
    try {
      setListLoadingStates(true)

      const { data } = await apis.listPartner1Goods(reqParams)
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

  useEffect(() => {
    fetchListGoods()
  }, [reqParams])

  const handleTableChange = (current: number, pageSize: number) => {
    setReqParams({ ...reqParams, currentPage: current, pageSize })
  };

  const syncGoods = async () => {
    setsyncLoadingStates(true)
    //@ts-ignore
    await apis.partner1SyncGoods({ id: ctx.cookie.id });
    message.success('同步成功');
    setsyncLoadingStates(false)
    onSearch({});
  }

  const onSearch = (value: IGoodsSearchCondition) => {
    setReqParams({ ...reqParams, ...value, currentPage: 1, pageSize: reqParams.pageSize })
  }

  return (
    <>
      <div style={{ marginBottom: '10px' }}>
        <CurrentLocation routeconfigs={getRouteConfig()} />
      </div>
      <Card>
        <div style={{ display: 'Flex' }}>
          <SearchForm onSearch={onSearch} />

          <Button type="primary" onClick={() => syncGoods()} loading={syncLoadingStates}>
            同步商品
          </Button>
          <></>
        </div>
        <div>
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
        </div>
      </Card>
    </>
  )
}

export default Goods

import { SearchOutlined } from "@ant-design/icons";
import { Button, Form, Select, SelectProps } from "antd";
import { useApis } from "../../api/api";
import { useEffect, useState } from "react";
import { IMerchant, ListMerchantReq } from "../../api/types";

export interface IMerchantSearchCondition {
  merchantId?: number;
}

interface SearchBarProps {
  onSearch: (conditions: ListMerchantReq) => void;
}

const MerchantSearchForm: React.FC<SearchBarProps> = ({ onSearch }) => {
  const [form] = Form.useForm<ListMerchantReq>();
  let { listMerchant } = useApis()

  const [idOptions, setIdOptions] = useState<SelectProps['options']>([]);

  const fetchListMerchant = async () => {
    const { data } = await listMerchant({ignoreStatistics: true})
    setMerchantIds(data.list)
  }

  const setMerchantIds = (datas: IMerchant[]) => {
    let ids: SelectProps['options'] = datas?.map((item) => {
      return {
        value: item.id,
        label: item.id + '(' + item.nickname + ')',
      }
    })

    setIdOptions(ids)
  }

  useEffect(() => {
    fetchListMerchant()
  }, [])

  const onFinish = (values: ListMerchantReq) => {
    onSearch(values);
  };

  return (
    <Form
      form={form}
      layout="inline"
      onFinish={onFinish}
      style={{ marginBottom: 16 }}
    >
      <Form.Item name="merchantId" label="商户ID">
        <Select
          allowClear
          showSearch
          size="middle"
          style={{ width: '200px' }}
          options={idOptions}
        />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit" icon={<SearchOutlined />}>
        </Button>
      </Form.Item>
    </Form>
  );
};

export default MerchantSearchForm
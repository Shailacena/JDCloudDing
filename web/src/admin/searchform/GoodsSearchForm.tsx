import { SearchOutlined } from "@ant-design/icons";
import { Button, Form, Select, SelectProps } from "antd";
import { useApis } from "../../api/api";
import { IPartner } from "../../api/types";
import { useEffect, useState } from "react";

export interface IGoodsSearchCondition {
  partnerId?: number;
}

interface SearchBarProps {
  onSearch: (conditions: IGoodsSearchCondition) => void;
}

const GoodsSearchForm: React.FC<SearchBarProps> = ({ onSearch }) => {
  const [form] = Form.useForm<IGoodsSearchCondition>();
  let { listPartner } = useApis()

  const [idOptions, setIdOptions] = useState<SelectProps['options']>([]);

  const fetchListPartner = async () => {
    const { data } = await listPartner({ ignoreStatistics: true })
    setPartnerIds(data.list)
  }

  const setPartnerIds = (datas: IPartner[]) => {
    let ids: SelectProps['options'] = datas?.map((item) => {
      return {
        value: item.id,
        label: item.id + '(' + item.nickname + ')',
      }
    })

    setIdOptions(ids)
  }

  useEffect(() => {
    fetchListPartner()
  }, [])

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
      <Form.Item name="partnerId" label="合作商ID">
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

export default GoodsSearchForm
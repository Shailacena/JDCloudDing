import { SearchOutlined } from "@ant-design/icons";
import { Button, DatePicker, Form, Select, SelectProps } from "antd";
import { useEffect, useState } from "react";
import { useApis } from "../../api/api";
import { IMerchant, ListDailyBillByMerchantReq } from "../../api/types";
import { Dayjs } from "dayjs";
import dayjs from "dayjs";
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';

// 扩展 dayjs 插件
dayjs.extend(utc);
dayjs.extend(timezone);

export interface IMerchantSearchCondition extends ListDailyBillByMerchantReq {
  merchantId?: number;
  dayRange?: Dayjs[];
}

interface SearchBarProps {
  onSearch: (conditions: ListDailyBillByMerchantReq) => void;
}

const MerchantDailyBillSearchForm: React.FC<SearchBarProps> = ({ onSearch }) => {
  const [form] = Form.useForm<IMerchantSearchCondition>();
  const [idOptions, setIdOptions] = useState<SelectProps['options']>([]);
  let { listMerchant } = useApis()

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

  const onFinish = (value: IMerchantSearchCondition) => {
    let { dayRange, ...params } = value
    onSearch({
      ...params,
      startAt: dayRange?.[0]?.tz('Asia/Shanghai').startOf('day').format('YYYY-MM-DD HH:mm:ss'),
      endAt: dayRange?.[1]?.tz('Asia/Shanghai').endOf('day').format('YYYY-MM-DD HH:mm:ss')
    });
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
      <Form.Item name="dayRange" label="日期">
        <DatePicker.RangePicker showNow={false} style={{ width: 250 }} />
      </Form.Item>
      <Form.Item>
        <Button type="primary" htmlType="submit" icon={<SearchOutlined />}>
        </Button>
      </Form.Item>
    </Form>
  );
};

export default MerchantDailyBillSearchForm
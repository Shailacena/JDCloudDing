import { SearchOutlined } from "@ant-design/icons";
import { Button, DatePicker, Form, Select, SelectProps } from "antd";
import { useEffect, useState } from "react";
import { useApis } from "../../api/api";
import { IPartner, ListDailyBillByPartnerReq } from "../../api/types";
import { Dayjs } from "dayjs";
import dayjs from "dayjs";
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';

// 扩展 dayjs 插件
dayjs.extend(utc);
dayjs.extend(timezone);

export interface IPantnerSearchCondition extends ListDailyBillByPartnerReq {
  partnerId?: number;
  dayRange?: Dayjs[];
}

interface SearchBarProps {
  onSearch: (conditions: ListDailyBillByPartnerReq) => void;
}

const PartnerDailyBillSearchForm: React.FC<SearchBarProps> = ({ onSearch }) => {
  const [form] = Form.useForm<IPantnerSearchCondition>();
  const [idOptions, setIdOptions] = useState<SelectProps['options']>([]);
  let { listPartner } = useApis()

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

  const onFinish = (value: IPantnerSearchCondition) => {
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
      <Form.Item name="partnerId" label="合作商ID">
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

export default PartnerDailyBillSearchForm
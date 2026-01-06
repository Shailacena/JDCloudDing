import { SearchOutlined } from "@ant-design/icons";
import { Button, DatePicker, Form } from "antd";
import { Dayjs } from "dayjs";
import dayjs from "dayjs";
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';

// 扩展 dayjs 插件
dayjs.extend(utc);
dayjs.extend(timezone);

interface SearchConf {
  OnSearch: Function
}

export interface ListBalanceBill {
  startAt?: string;
  endAt?: string;
}

export interface IRecordSearchCondition extends ListBalanceBill {
  dayRange: Dayjs[]
}

const BalanceBillSearchForm = (searchConf: SearchConf) => {
  const onFinish = (value: IRecordSearchCondition) => {
    let { dayRange, ...params } = value
    searchConf.OnSearch?.({
      ...params,
      startAt: dayRange?.[0]?.tz('Asia/Shanghai').startOf('day').format('YYYY-MM-DD HH:mm:ss'),
      endAt: dayRange?.[1]?.tz('Asia/Shanghai').endOf('day').format('YYYY-MM-DD HH:mm:ss')
    })
  };

  return (
    <Form
      layout="inline"
      onFinish={onFinish}
      style={{ marginBottom: 16 }}
    >

      <Form.Item<IRecordSearchCondition> name="dayRange">
        <DatePicker.RangePicker style={{ width: 250 }} />
      </Form.Item>

      <Form.Item>
        <Button type="primary" htmlType="submit" icon={<SearchOutlined />}></Button>
      </Form.Item>
    </Form>
  );
};

export default BalanceBillSearchForm
import { SearchOutlined } from "@ant-design/icons";
import { Button, DatePicker, Form, Input } from "antd";
import { ListOrderReq } from "../../api/types";
import { Dayjs } from "dayjs";
import dayjs from "dayjs";
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';

// 扩展 dayjs 插件
dayjs.extend(utc);
dayjs.extend(timezone);

interface ItemConf {
  hiddenPart?: boolean
  OnSearch?: Function
}

export interface IRecordSearchCondition extends ListOrderReq {
  dayRange: Dayjs[]
}

const TradingRecordSearchForm = (itemConf: ItemConf) => {
  const onFinish = (value: IRecordSearchCondition) => {
    let { dayRange, ...params } = value
    itemConf.OnSearch?.({
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
      <Form.Item<IRecordSearchCondition> name="partnerOrderId">
        <Input placeholder="店铺单号" allowClear style={{ width: 150, marginLeft: 0, marginRight: 0 }} />
      </Form.Item>
      <Form.Item<IRecordSearchCondition> name="orderId">
        <Input placeholder="系统单号" allowClear style={{ width: 150, marginLeft: 0, marginRight: 0 }} />
      </Form.Item>

      {
        !itemConf.hiddenPart && <Form.Item<IRecordSearchCondition> name="dayRange">
          <DatePicker.RangePicker style={{ width: 250 }} />
        </Form.Item>
      }
      <Form.Item>
        <Button type="primary" htmlType="submit" icon={<SearchOutlined />}></Button>
      </Form.Item>
    </Form>
  );
};

export default TradingRecordSearchForm

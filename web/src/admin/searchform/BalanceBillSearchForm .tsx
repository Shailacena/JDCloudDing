import { SearchOutlined } from "@ant-design/icons";
import { Button, DatePicker, Form, Select, SelectProps } from "antd";
import { Dayjs } from "dayjs";
import dayjs from "dayjs";
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';

// 扩展 dayjs 插件
dayjs.extend(utc);
dayjs.extend(timezone);

export enum SearchType {
  Partner,
  Merchant
}

export interface SearchFormIdItem {
  id: number
  label: string
}

interface SearchConf {
  searchType: SearchType
  OnSearch: (req: ListBalanceBill) => void;
  ids: SearchFormIdItem[]
}

export interface ListBalanceBill {
  id?: number;
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

  const idOptions: SelectProps['options'] = [];

  searchConf.ids.forEach((item) => {
    idOptions.push({
      value: item.id,
      label: item.label
    })
  })

  return (
    <Form
      layout="inline"
      onFinish={onFinish}
      style={{ marginBottom: 16 }}
    >

      <Form.Item<IRecordSearchCondition>
        name="id"
        label={searchConf.searchType === SearchType.Partner ? '合作商ID' : '商户ID'}
      >
        <Select
          allowClear
          showSearch
          size="middle"
          style={{ width: '200px' }}
          options={idOptions}
        />
      </Form.Item>
      {
        <Form.Item<IRecordSearchCondition> name="dayRange">
          <DatePicker.RangePicker style={{ width: 250 }} />
        </Form.Item>
      }
      <Form.Item>
        <Button type="primary" htmlType="submit" icon={<SearchOutlined />}></Button>
      </Form.Item>
    </Form>
  );
};

export default BalanceBillSearchForm
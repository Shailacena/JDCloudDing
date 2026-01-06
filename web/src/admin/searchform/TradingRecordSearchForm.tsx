import { CloseCircleOutlined, SearchOutlined, DownloadOutlined } from "@ant-design/icons";
import { Button, DatePicker, Form, Input, Select, SelectProps, Space, message } from "antd";
import { IMerchant, IPartner, ListOrderReq } from "../../api/types";
import { useEffect, useState } from "react";
import { useApis } from "../../api/api";
import { Dayjs } from "dayjs";
import dayjs from "dayjs";
import utc from 'dayjs/plugin/utc';
import timezone from 'dayjs/plugin/timezone';

// 扩展 dayjs 插件
dayjs.extend(utc);
dayjs.extend(timezone);

interface ItemConf {
  hiddenPart?: boolean
  OnSearch?: (req: ListOrderReq) => void;
  OnDownload?: (req: ListOrderReq) => void;
}

export interface IRecordSearchCondition extends ListOrderReq {
  dayRange: Dayjs[]
}

const TradingRecordSearchForm = (itemConf: ItemConf) => {
  const [partnerIdsOptions, setPartnerIdsOptions] = useState<SelectProps['options']>([]);
  const [merchantIdsOptions, setMerchantIdsOptions] = useState<SelectProps['options']>([]);
  const [form] = Form.useForm();

  let { listPartner, listMerchant } = useApis()

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

    setPartnerIdsOptions(ids)
  }

  const fetchListMerchant = async () => {
    const { data } = await listMerchant({ ignoreStatistics: true })
    setMerchantIds(data.list)
  }

  const setMerchantIds = (datas: IMerchant[]) => {
    let ids: SelectProps['options'] = datas?.map((item) => {
      return {
        value: item.id,
        label: item.id + '(' + item.nickname + ')',
      }
    })

    setMerchantIdsOptions(ids)
  }

  useEffect(() => {
    fetchListPartner()
    fetchListMerchant()
  }, [])

  const onFinish = (value: IRecordSearchCondition) => {
    let { dayRange, ...params } = value
    itemConf.OnSearch?.({
      ...params,
      startAt: dayRange?.[0]?.tz('Asia/Shanghai').startOf('day').format('YYYY-MM-DD HH:mm:ss'),
      endAt: dayRange?.[1]?.tz('Asia/Shanghai').endOf('day').format('YYYY-MM-DD HH:mm:ss')
    })
  };

  const resetFormFields = () => {
    form.resetFields()
    itemConf.OnSearch?.({} as ListOrderReq)
  }

  const onDownload = () => {
    const formValues = form.getFieldsValue() as IRecordSearchCondition;
    const { dayRange } = formValues;

    // 验证日期选择
    if (!dayRange || dayRange.length === 0) {
      // 如果没有选择日期，使用当前日期
      const today = dayjs().tz('Asia/Shanghai');
      itemConf.OnDownload?.({
        ...formValues,
        startAt: today.startOf('day').format('YYYY-MM-DD HH:mm:ss'),
        endAt: today.endOf('day').format('YYYY-MM-DD HH:mm:ss')
      });
      return;
    }

    if (dayRange.length === 1) {
      // 只选择了一个日期，使用这个日期
      const selectedDate = dayRange[0].tz('Asia/Shanghai');
      itemConf.OnDownload?.({
        ...formValues,
        startAt: selectedDate.startOf('day').format('YYYY-MM-DD HH:mm:ss'),
        endAt: selectedDate.endOf('day').format('YYYY-MM-DD HH:mm:ss')
      });
      return;
    }

    if (dayRange.length === 2) {
      // 选择了日期范围，检查是否为同一天
      const startDate = dayRange[0].tz('Asia/Shanghai').format('YYYY-MM-DD');
      const endDate = dayRange[1].tz('Asia/Shanghai').format('YYYY-MM-DD');
      
      if (startDate === endDate) {
        // 同一天，允许下载
        itemConf.OnDownload?.({
          ...formValues,
          startAt: dayRange[0].tz('Asia/Shanghai').startOf('day').format('YYYY-MM-DD HH:mm:ss'),
          endAt: dayRange[1].tz('Asia/Shanghai').endOf('day').format('YYYY-MM-DD HH:mm:ss')
        });
      } else {
        // 不同日期，提示错误
        message.error('下载功能只支持单日数据导出，请选择同一天的日期范围');
      }
      return;
    }

    message.error('请选择有效的日期');
  }

  return (
    <Form
      form={form}
      layout="inline"
      onFinish={onFinish}
      style={{ marginBottom: 16 }}
    >
      <Form.Item<IRecordSearchCondition> name="orderId">
        <Input placeholder="系统单号" allowClear style={{ width: 150, marginLeft: 0, marginRight: 0 }} />
      </Form.Item>
      <Form.Item<IRecordSearchCondition> name="partnerOrderId">
        <Input placeholder="店铺单号" allowClear style={{ width: 150, marginLeft: 0, marginRight: 0 }} />
      </Form.Item>
      <Form.Item<IRecordSearchCondition> name="merchantOrderId">
        <Input placeholder="商户单号" allowClear style={{ width: 150 }} />
      </Form.Item>
      <Form.Item<IRecordSearchCondition> name="partnerId">
        <Select
          allowClear
          showSearch
          size="middle"
          placeholder="合作商ID"
          style={{ width: '200px' }}
          options={partnerIdsOptions}
        />
      </Form.Item>
      <Form.Item<IRecordSearchCondition> name="merchantId">
        <Select
          allowClear
          showSearch
          size="middle"
          placeholder="商户ID"
          style={{ width: '200px' }}
          options={merchantIdsOptions}
        />
      </Form.Item>
      {
        !itemConf.hiddenPart && <Form.Item<IRecordSearchCondition> name="dayRange">
          <DatePicker.RangePicker style={{ width: 250 }} />
        </Form.Item>
      }
      <Form.Item>
        <Space size="small">
          <Button type="primary" htmlType="submit" icon={<SearchOutlined />}></Button>
          <Button type="primary" onClick={resetFormFields} icon={<CloseCircleOutlined />}></Button>
          {itemConf.OnDownload && (
            <Button type="default" onClick={onDownload} icon={<DownloadOutlined />}>
              下载
            </Button>
          )}
        </Space>
      </Form.Item>
    </Form>
  );
};

export default TradingRecordSearchForm

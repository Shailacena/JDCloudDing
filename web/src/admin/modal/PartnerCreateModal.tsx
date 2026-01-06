import { useEffect, useState } from 'react';
import { Button, Form, FormProps, Input, message, Modal, Select, Row, Col } from 'antd';
import { useApis } from '../../api/api';
import { IPartner, PartnerBaseInfoReq, PartnerUpdateReq } from '../../api/types';
import axios from 'axios';
import TextArea from 'antd/es/input/TextArea';
import { useAppContext } from '../../AppProvider';
import { AllChannelId, AllPartnerType, ChannelPayType, PartnerType } from '../../utils/biz';

// 定义选项类型
interface OptionType {
  value: string | number;
  label: string;
}

// 定义Modal数据类型
interface ModalDataType {
  isModalOpen: boolean;
  info?: IPartner;
  onOk?: () => void;
  onCancel?: () => void;
}

// 定义Title常量
const Title = {
  EditTxt: '编辑合作商',
  CreateTxt: '添加合作商'
};

// 判断是否为京东复制通道
const isJDPayChannel = (channelId?: string): boolean => channelId === 'JS000000';

// 获取可用的合作商类型选项
const getAvailablePartnerTypes = (channelId?: string): OptionType[] => {
  // 京东复制通道只能选择阿奇索类型
  if (isJDPayChannel(channelId) && AllPartnerType) {
    return AllPartnerType
      .filter(item => item.id === PartnerType.Agiso)
      .map(item => ({
        value: item.id,
        label: item.label
      }));
  }
  // 其他通道可以选择阿奇索或安式
  return AllPartnerType ? AllPartnerType.map(item => ({
    value: item.id,
    label: item.label
  })) : [];
};

// 获取可用的跳转类型选项
const getAvailablePayTypes = (channelId?: string): OptionType[] => {
  // 确保channelId存在且ChannelPayType[channelId]有效
  if (!channelId || !ChannelPayType || !ChannelPayType[channelId]) {
    return [];
  }
  
  // 京东复制通道只有app跳转选项
  if (isJDPayChannel(channelId)) {
    return ChannelPayType[channelId]
      .filter((item: any) => item.payType === 'app')
      .map((item: any) => ({
        value: item.payType,
        label: item.label
      }));
  }
  // 其他通道保留原有选项
  return ChannelPayType[channelId]
    .map((item: any) => ({
      value: item.payType,
      label: item.label
    }));
};

const PartnerCreateModal = (params: ModalDataType) => {
  const [info, setInfo] = useState<IPartner | undefined>(params.info)
  const [isEdit, setIsEdit] = useState<boolean>(!!params.info)
  const [title, setTitle] = useState<string>('')
  const [isModalOpen, setIsModalOpen] = useState<boolean>(params.isModalOpen);
  const [confirmLoading, setConfirmLoading] = useState<boolean>(false);
  const [formDisabled, setFormDisabled] = useState<boolean>(false);
  const [jumpOptions, setJumpOptions] = useState<OptionType[]>([]);
  const [partnerTypeOptions, setPartnerTypeOptions] = useState<OptionType[]>([]);
  const [selectedChannelId, setSelectedChannelId] = useState<string | undefined>(params.info?.channelId);
  const app = useAppContext()
  const apis = useApis()
  const [form] = Form.useForm<PartnerBaseInfoReq & PartnerUpdateReq>();
  const [partnerType, setPartnerType] = useState<number | undefined>(params.info?.type);

  const handleChannelChange = (value: string): void => {
    setSelectedChannelId(value);
    
    // 如果是京东复制通道，直接设置默认选项
    if (isJDPayChannel(value)) {
      // 直接设置阿奇索类型选项
      const agisoOption = { value: PartnerType.Agiso, label: '阿奇索' };
      setPartnerTypeOptions([agisoOption]);
      setPartnerType(PartnerType.Agiso);
      
      // 直接设置App跳转选项
      setJumpOptions([{ value: 1, label: 'App跳转' }]); // 使用数字类型匹配接口定义
      
      // 主动设置表单值，确保类型选择正确显示
      form.setFieldsValue({
        type: PartnerType.Agiso,
        payType: 1 // 使用数字类型匹配接口定义
      });
    } else {
      // 更新可用的跳转类型选项
      setJumpOptions(getAvailablePayTypes(value));
      // 更新可用的合作商类型选项
      const types = getAvailablePartnerTypes(value);
      setPartnerTypeOptions(types);
      
      // 清除之前设置的类型值，让用户重新选择
      form.setFieldsValue({
        type: undefined,
        payType: undefined
      });
    }
  };

  useEffect(() => {
    setIsModalOpen(params.isModalOpen)
  }, [params.isModalOpen])

  useEffect(() => {
    setIsEdit(!!params.info)
    setInfo(params.info)
    setSelectedChannelId(params.info?.channelId)
    setPartnerType(params.info?.type)
    
    // 初始化选项
    if (params.info?.channelId) {
      // 如果是京东复制通道，直接设置默认选项
      if (isJDPayChannel(params.info.channelId)) {
        // 直接设置阿奇索类型选项
        const agisoOption = { value: PartnerType.Agiso, label: '阿奇索' };
        setPartnerTypeOptions([agisoOption]);
        setPartnerType(PartnerType.Agiso);
        
        // 直接设置App跳转选项
        setJumpOptions([{ value: 1, label: 'App跳转' }]); // 使用数字类型匹配接口定义
      } else {
        // 其他通道使用原来的逻辑
        setJumpOptions(getAvailablePayTypes(params.info.channelId));
        setPartnerTypeOptions(getAvailablePartnerTypes(params.info.channelId));
      }
    }
  }, [params.info])

  useEffect(() => {
    setTitle(isEdit ? Title.EditTxt : Title.CreateTxt)
    // 初始化类型选项
    setPartnerTypeOptions(getAvailablePartnerTypes(selectedChannelId));
  }, [isEdit, selectedChannelId])

  const onFinish: FormProps<PartnerBaseInfoReq & PartnerUpdateReq>['onFinish'] = async (value) => {
    setFormDisabled(true)
    setConfirmLoading(true)
    try {
      if (isEdit && params.info?.id) {
        // 先展开value，再设置id以避免重复定义警告
        const updateValue = { ...value, id: params.info.id } as PartnerUpdateReq;
        handleEdit(updateValue);
      } else {
        handleRegister(value as PartnerBaseInfoReq);
      }
    } catch (e) {
      if (axios.isAxiosError(e)) {
        let msg = e.response?.data?.message;
        if (msg) message.error(msg);
      }
    } finally {
      setFormDisabled(false)
      setConfirmLoading(false)
    }
  };

  const handleRegister = async (value: PartnerBaseInfoReq): Promise<void> => {
    value.priority = value.priority ? Number(value.priority) : 0;
    value.level = app.cookie?.level || 1;
    const { data } = await apis.partnerRegister(value);
    params.onOk?.();
    Modal.success({
      content: `添加成功, 密位为${data.password}`,
    });
    if (app.partnerList) {
      app.partnerList = [];
    }
  };

  const handleEdit = async (value: PartnerUpdateReq): Promise<void> => {
    if (value.priority !== undefined) {
      value.priority = Number(value.priority);
    }
    if (value.rechargeTime !== undefined) {
      value.rechargeTime = Number(value.rechargeTime);
    }
    await apis.partnerUpdate(value);
    params.onOk?.();
    message.success(`修改成功`);
  };

  const isType = (value?: number, type?: number): boolean => {
    return value === type;
  };

  const handlePartnerTypeChange = (value: number): void => {
    setPartnerType(value);
  };

  return (
    <>
      <Modal 
        title={title}
        okText="确认"
        cancelText="取消"
        onCancel={() => {
          form.resetFields();
          params.onCancel?.();
        }}
        confirmLoading={confirmLoading}
        open={isModalOpen}
        // 调整宽度以适应两列布局
        width={800}
        style={{ maxWidth: '95%' }}
        footer={[
          <Button key="cancel" onClick={() => {
            form.resetFields();
            params.onCancel?.();
          }}>取消</Button>,
          <Button
            key="submit"
            type="primary"
            loading={confirmLoading}
            onClick={() => {
              form.submit();
            }}
          >
            确认
          </Button>,
        ]}
      >
        <Form
          form={form}
          layout="vertical"
          disabled={formDisabled}
          onFinish={onFinish}
          initialValues={{
            ...info,
            priority: info?.priority ?? 10,
            darkNumberLength: info?.darkNumberLength ?? 11
          }}
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item<PartnerBaseInfoReq & PartnerUpdateReq>
                name="nickname"
                label="名称"
                rules={[{ required: true, message: '请输入名称' }]}
              >
                <Input />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item<PartnerBaseInfoReq & PartnerUpdateReq>
                name="channelId"
                label="通道ID"
                rules={[{ required: !isEdit, message: '请选择通道ID' }]}
              >
                <Select
                  disabled={isEdit}
                  options={
                    AllChannelId?.map((item) => ({
                      value: item.channelId,
                      label: item.label
                    })) || []
                  }
                  onChange={handleChannelChange}
                />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item<PartnerBaseInfoReq & PartnerUpdateReq>
                name="type"
                label="类型"
                rules={[{ required: true, message: '请选择类型' }]}
                initialValue={isJDPayChannel(selectedChannelId) ? PartnerType.Agiso : undefined}
              >
                <Select
                  disabled={isEdit} // 只在编辑模式下禁用，京东复制通道下允许点击
                  options={partnerTypeOptions}
                  onChange={handlePartnerTypeChange}
                  // 确保京东复制通道下始终显示阿奇索选项
                  value={isJDPayChannel(selectedChannelId) ? PartnerType.Agiso : undefined}
                />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item<PartnerBaseInfoReq & PartnerUpdateReq>
                name="payType"
                label="跳转类型"
                rules={[{ required: !isEdit, message: '请选择跳转类型' }]}
              >
                <Select options={jumpOptions} />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item<PartnerBaseInfoReq & PartnerUpdateReq>
                name="priority"
                label="优先级"
                initialValue={10}
                rules={[{ required: true, message: '请输入优先级' }]}
              >
                <Input type="number" min={0} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item<PartnerBaseInfoReq & PartnerUpdateReq>
                name="darkNumberLength"
                label="帐号位数"
                initialValue={11}
                rules={[
                  { required: true, message: '请输入帐号位数' },
                  {
                    validator: (_, value) => {
                      const num = Number(value);
                      if (isNaN(num)) {
                        return Promise.reject('帐号位数必须是数字');
                      }
                      if (num < 8 || num > 15) {
                        return Promise.reject('帐号位数必须在8-15之间');
                      }
                      return Promise.resolve();
                    }
                  }
                ]}
                getValueFromEvent={(e) => {
                  const value = e.target.value;
                  return value ? Number(value) : undefined;
                }}
              >
                <Input type="number" min={8} max={15} placeholder="8-15位，默认11位" />
              </Form.Item>
            </Col>
          </Row>

          {/* 阿奇索配置 - 两列布局 */}
          {isType(partnerType, PartnerType.Agiso) && (
            <Row gutter={16}>
              <Col span={12}>
                <Form.Item<PartnerBaseInfoReq & PartnerUpdateReq>
                  name="aqsAppSecret"
                  label="阿奇索Secret"
                  rules={[{ required: true, message: '请输入阿奇索Secret' }]}
                >
                  <Input />
                </Form.Item>
              </Col>
              <Col span={12}>
                <Form.Item<PartnerBaseInfoReq & PartnerUpdateReq>
                  name="aqsToken"
                  label="阿奇索Token"
                  rules={[{ required: true, message: '请输入阿奇索Token' }]}
                >
                  <Input />
                </Form.Item>
              </Col>
            </Row>
          )}

          {/* 安式配置 - 两列布局 */}
          {isType(partnerType, PartnerType.Anssy) && isEdit && (
            <Row gutter={16}>
              <Col span={12}>
                <Form.Item<PartnerBaseInfoReq & PartnerUpdateReq>
                  name="anssyAppSecret"
                  label="安式Secret"
                >
                  <Input disabled={isEdit} />
                </Form.Item>
              </Col>
              <Col span={12}>
                <Form.Item<PartnerBaseInfoReq & PartnerUpdateReq>
                  name="anssyToken"
                  label="安式Token"
                >
                  <Input disabled={isEdit} />
                </Form.Item>
              </Col>
            </Row>
          )}

          {/* 私钥 - 单独一行 */}
          {isEdit && (
            <Form.Item<PartnerBaseInfoReq & PartnerUpdateReq>
              name="privateKey"
              label="私钥"
            >
              <Input disabled />
            </Form.Item>
          )}
          
          {/* 备注 - 单独一行 */}
          <Form.Item<PartnerBaseInfoReq & PartnerUpdateReq>
            name="remark"
            label="备注"
          >
            <TextArea rows={4} />
          </Form.Item>
        </Form>
      </Modal>
    </>
  );
}

export default PartnerCreateModal
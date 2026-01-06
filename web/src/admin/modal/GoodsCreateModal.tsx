import { useEffect, useState } from 'react';
import { Button, Form, FormProps, Input, message, Modal, Select, Switch } from 'antd';
import { GoodsCreateReq, GoodsUpdateReq } from '../../api/types';
import { useApis } from '../../api/api';
import axios from 'axios';
import { useAppContext } from '../../AppProvider';
import { GoodsStatus } from '../../utils/biz';

interface ModalDataType {
  isModalOpen: boolean
  onOk: Function;
  onCancel: Function;
  info?: FieldType
}

export type FieldType = {
  id?: number;
  partnerId?: number;
  skuId?: string;
  price?: number;
  realPrice?: number;
  // shopName?: string;
  status?: number;
};

enum Title {
  CreateTxt = '新增商品',
  EditTxt = '修改商品'
}

const GoodsCreateModal = (params: ModalDataType) => {
  const [info, setInfo] = useState(params.info)
  const [isEdit, setIsEdit] = useState(!!params.info)
  const [title, setTitle] = useState('')
  const [isModalOpen, setIsModalOpen] = useState(params.isModalOpen);
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [formDisabled, setFormDisabled] = useState<boolean>(false);
  // 使用组件内部状态管理partnerList，而不是依赖全局context
  const [partnerList, setPartnerList] = useState<any[]>([]);
  let { listPartner, createGoods, goodsUpdate } = useApis()
  const ctx = useAppContext();
  const [form] = Form.useForm();
  useEffect(() => {
    setIsModalOpen(params.isModalOpen)
  }, [params.isModalOpen])

  useEffect(() => {
    setIsEdit(!!params.info)
    setInfo(params.info)
  }, [params.info])

  useEffect(() => {
    setTitle(isEdit ? Title.EditTxt : Title.CreateTxt)
  }, [isEdit])

  useEffect(() => {
    fetchListPartner()
  }, [])

  const fetchListPartner = async () => {
    try {
      const { data } = await listPartner({ ignoreStatistics: true })
      // 更新组件内部状态
      setPartnerList(data?.list || []);
      // 如果需要，仍然可以尝试更新全局context，但这不是主要的状态来源
      if (ctx && data?.list) {
        try {
          ctx.partnerList = data.list;
        } catch (e) {
          console.warn('无法更新全局partnerList状态');
        }
      }
    } catch (error) {
      console.error('获取合作商列表失败:', error);
      message.error('获取合作商列表失败');
    }
  }

  const onFinish: FormProps<GoodsUpdateReq>['onFinish'] = async (value) => {
    setFormDisabled(true)
    setConfirmLoading(true)
    value.price = +value.price
    // 修改realPrice逻辑，当realPrice为空时使用price值
    if (!value.realPrice && value.realPrice !== 0) {
      value.realPrice = value.price;
    } else {
      value.realPrice = +value.realPrice;
    }
    // 修正status值的处理，确保与GoodsStatus枚举对应
    value.status = value.status ? GoodsStatus.Enabled : GoodsStatus.Disabled
    try {
      isEdit ? handleEdit(value) : handleRegister(value)
    } catch (e) {
      if (axios.isAxiosError(e)) {
        let msg = e.response?.data?.message
        msg && message.error(msg)
      }
    } finally {
      setFormDisabled(false) 
      setConfirmLoading(false)
    }
  };

  const handleRegister = async (value: GoodsCreateReq) => {
    await createGoods(value)
    params?.onOk?.();
    message.success(`添加成功`)
  }

  const handleEdit = async (value: GoodsUpdateReq) => {
    if (info?.id) {
      value.id = info?.id
    }
    await goodsUpdate({ ...value })
    params?.onOk?.();
    message.success(`修改成功`)
  }

  return (
    <Modal
      open={isModalOpen}
      title={title}
      okText="确认"
      cancelText="取消"
      onCancel={() => {
        form.resetFields();
        params.onCancel();
      }}
      confirmLoading={confirmLoading}
      // 添加宽度限制，减小模态框宽度
      width={500}
      style={{ maxWidth: '90%' }}
      footer={[
        <Button key="cancel" onClick={() => {
          form.resetFields();
          params.onCancel();
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
        initialValues={{
          // 合并info中的其他初始值
          ...info,
          // 特殊处理status字段，根据info中的status值设置正确的开关状态
          status: info?.status === GoodsStatus.Enabled || info?.status === undefined ? true : false
        }}
        onFinish={onFinish}
        disabled={formDisabled}
      >
        <Form.Item
          label="合作商"
          name="partnerId"
          rules={[{ required: true, message: '请选择合作商' }]}
        >
          <Select placeholder="请选择合作商">
            {partnerList?.map((item) => (
              <Select.Option key={item.id} value={item.id}>
                {item.id + `(${item.nickname})`} 
              </Select.Option>
            ))}
          </Select>
        </Form.Item>
        <Form.Item
          label="商品编号（SKU）"
          name="skuId"
          rules={[{ required: true, message: '请输入商品编号' }]}
        >
          <Input placeholder="请输入商品编号" />
        </Form.Item>
        <Form.Item
          label="金额"
          name="price"
          rules={[{ required: true, message: '请输入金额' }]}
        >
          <Input
            placeholder="请输入金额"
            type="number"
          />
        </Form.Item>
        <Form.Item
          label="真实金额"
          name="realPrice"
          rules={[{ required: false }]}
        >
          <Input placeholder="请输入真实金额" type="number" />
        </Form.Item>
        <Form.Item
          label="状态"
          name="status"
          valuePropName="checked"
        >
          <Switch />
        </Form.Item>
      </Form>
    </Modal>
  )
}

export default GoodsCreateModal

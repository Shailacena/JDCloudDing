import { Button, Form, FormProps, Input, message, Modal, Typography } from "antd";
import { useApis } from "../../api/api";
import { UpdateMerchantBalanceReq } from "../../api/types";
import { useState } from "react";
import axios from "axios";
import { useAppContext } from "../../AppProvider";

interface ModalDataType {
  isModalOpen: boolean
  onOk: Function;
  onCancel: Function;
  merchantId: number
}

const UpdateMerchantBalanceModal = (params: ModalDataType) => {
  let { merchantUpdateBalance } = useApis()
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [formDisabled, setFormDisabled] = useState<boolean>(false);
  let app = useAppContext()
  const [form] = Form.useForm();

  const onFinish: FormProps<UpdateMerchantBalanceReq>['onFinish'] = async (value) => {
    setFormDisabled(true)
    setConfirmLoading(true)

    try {
      updateBalance(value)
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

  const updateBalance = async (value: UpdateMerchantBalanceReq) => {
    value.adminId = app.cookie.id || 0
    value.merchantId = params.merchantId
    value.changeAmount = +value.changeAmount
    await merchantUpdateBalance(value)

    params?.onOk?.();
    message.success(`调整余额成功`)
  }

  return (
    <>
      <Modal
        title="余额调整"
        okText="确认"
        cancelText="取消"
        onCancel={() => {
          form.resetFields();
          params?.onCancel?.();
        }}
        confirmLoading={confirmLoading}
        open={params.isModalOpen}
        // 添加宽度限制
        width={500}
        style={{ maxWidth: '90%' }}
        footer={[
          <Button key="cancel" onClick={() => {
            form.resetFields();
            params?.onCancel?.();
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
        >
          <Form.Item<UpdateMerchantBalanceReq>
            name="changeAmount"
            label="调整金额"
            rules={[{ required: true, message: '请输入金额' }]}
            help={<Typography.Text type="danger">整数增加，负数减少</Typography.Text>}
          >
            <Input type="number" />
          </Form.Item>

          <Form.Item<UpdateMerchantBalanceReq>
            name="password"
            label="登录密码"
            rules={[{ required: true, message: '请输入登录密码' }]}
          >
            <Input.Password />
          </Form.Item>
        </Form >
      </Modal>
    </>
  );
};

export default UpdateMerchantBalanceModal;
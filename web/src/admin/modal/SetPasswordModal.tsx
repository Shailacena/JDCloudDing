import { Button, Form, FormProps, Input, message, Modal } from "antd";
import { useApis } from "../../api/api";
import { AdminSetPasswordReq } from "../../api/types";
import { useState } from "react";
import axios from "axios";

interface ModalDataType {
  isModalOpen: boolean
  onOk: Function;
  onCancel: Function;
}

const SetPasswordModal = (params: ModalDataType) => {
  let { adminSetPassword } = useApis()
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [formDisabled, setFormDisabled] = useState<boolean>(false);
  const [form] = Form.useForm();

  const onFinish: FormProps<AdminSetPasswordReq>['onFinish'] = async (value) => {
    setFormDisabled(true)
    setConfirmLoading(true)

    try {
      handleSetPassword(value)
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

  const handleSetPassword = async (value: AdminSetPasswordReq) => {
    await adminSetPassword(value)

    params?.onOk?.();
    message.success(`修改密码成功`)
  }

  return (
    <>
      <Modal
        title="修改密码"
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
          <Form.Item<AdminSetPasswordReq>
            name="oldPassword"
            label="原密码"
            rules={[{ required: true, message: '请输入原密码' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item<AdminSetPasswordReq>
            name="newPassword"
            label="新密码"
            rules={[{ required: true, message: '请输入新密码' }]}
          >
            <Input />
          </Form.Item>
        </Form >
      </Modal>
    </>
  );
};

export default SetPasswordModal;
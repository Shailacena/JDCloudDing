import { useState } from 'react';
import { Button, Divider, Flex, Form, FormProps, Input, message, Modal } from 'antd';
import { useApis } from '../../api/api';
import axios from 'axios';
import { MerchantSetPasswordReq } from '../../api/types';

interface MerchantSetPasswordDataType {
  isModalOpen: boolean
  onOk: Function;
  onCancel: Function;
}

const SetPasswordModal = (params: MerchantSetPasswordDataType) => {
  let { merchant1SetPassword: merchantSetPassword } = useApis()
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [formDisabled, setFormDisabled] = useState<boolean>(false);

  const onFinish: FormProps<MerchantSetPasswordReq>['onFinish'] = async (value) => {
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

  const handleSetPassword = async (value: MerchantSetPasswordReq) => {
    await merchantSetPassword(value)

    params?.onOk?.();
    message.success(`修改密码成功`)
  }

  return (
    <>
      <Modal
        title="修改密码"
        footer={null}
        confirmLoading={confirmLoading}
        open={params.isModalOpen}
        onCancel={() => { params?.onCancel?.() }}
        destroyOnClose
      >
        <Divider />
        <Form
          labelCol={{ span: 4 }}
          name="basic"
          autoComplete="off"
          disabled={formDisabled}
          onFinish={onFinish}
        >
          <Form.Item<MerchantSetPasswordReq>
            name="oldPassword"
            label="原密码"
            rules={[{ required: true, message: '请输入原密码' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item<MerchantSetPasswordReq>
            name="newPassword"
            label="新密码"
            rules={[{ required: true, message: '请输入新密码' }]}
          >
            <Input />
          </Form.Item>

          <Form.Item>
            <Flex justify="center" align="center">
              <Button size="large" type="primary" htmlType="submit">
                确定
              </Button>
            </Flex>
          </Form.Item>
        </Form >
      </Modal>
    </>
  );
};

export default SetPasswordModal;
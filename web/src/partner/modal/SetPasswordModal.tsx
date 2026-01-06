import { useState } from 'react';
import { Button, Divider, Flex, Form, FormProps, Input, message, Modal } from 'antd';
import { useApis } from '../../api/api';
import axios from 'axios';
import { PartnerSetPasswordReq } from '../../api/types';

interface PartnerSetPasswordDataType {
  isModalOpen: boolean
  onOk: Function;
  onCancel: Function;
}

const SetPasswordModal = (params: PartnerSetPasswordDataType) => {
  let { partner1SetPassword: partnerSetPassword } = useApis()
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [formDisabled, setFormDisabled] = useState<boolean>(false);

  const onFinish: FormProps<PartnerSetPasswordReq>['onFinish'] = async (value) => {
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

  const handleSetPassword = async (value: PartnerSetPasswordReq) => {
    await partnerSetPassword(value)

    params?.onOk?.();
    message.success(`修改密码成功`)
  }

  return (
    <>
      <Modal title="修改密码"
        footer={null}
        open={params.isModalOpen}
        confirmLoading={confirmLoading}
        onCancel={() => { params?.onCancel?.() }}
        style={{ maxWidth: 480 }}
        destroyOnClose>
        <Divider />
        <div style={{ display: 'flex', justifyContent: 'center', alignContent: 'center', marginTop: 20, alignItems: 'center' }}>
          <Form
            labelCol={{ span: 8 }}
            name="basic"
            autoComplete="off"
            onFinish={onFinish}
            disabled={formDisabled}
          >
            <Form.Item<PartnerSetPasswordReq>
              name="oldpassword"
              label="原始密码"
              required
            >
              <Input style={{ width: 200 }} />
            </Form.Item>

            <Form.Item<PartnerSetPasswordReq>
              name="newpassword"
              label="新密码"
              required
            >
              <Input style={{ width: 200 }} />
            </Form.Item>

            <Form.Item>
              <Flex justify="center" align="center">
                <Button size="large" type="primary" htmlType="submit">
                  确定
                </Button>
              </Flex>
            </Form.Item>
          </Form >
        </div>
      </Modal>
    </>
  );
};

export default SetPasswordModal;
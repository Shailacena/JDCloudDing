import { useEffect, useState } from 'react';
import { Button, Form, FormProps, Input, message, Modal } from 'antd';
import { MerchantRegisterReq, MerchantUpdateReq } from '../../api/types';
import { useApis } from '../../api/api';
import axios from 'axios';
import { useAppContext } from '../../AppProvider';

const { TextArea } = Input;

interface ModalDataType {
  isModalOpen: boolean
  onOk: Function;
  onCancel: Function;
  info?: FieldType
}

export type FieldType = {
  id?: number;
  username?: string;
  nickname?: string;
  privateKey?: string;
  createAt?: number;
  totalAmount?: number;
  todayAmount?: number;
  enable?: number;
  remark?: string;
};

enum Title {
  CreateTxt = '新增商户',
  EditTxt = '修改商户'
}

const MerchantCreateModal = (params: ModalDataType) => {
  const [info, setInfo] = useState(params.info)
  const [isEdit, setIsEdit] = useState(!!params.info)
  const [title, setTitle] = useState('')
  const [isModalOpen, setIsModalOpen] = useState(params.isModalOpen);
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [formDisabled, setFormDisabled] = useState<boolean>(false);
  let { merchantRegister, merchantUpdate } = useApis()
  let app = useAppContext()
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

  const onFinish: FormProps<MerchantRegisterReq>['onFinish'] = async (value) => {
    setFormDisabled(true)
    setConfirmLoading(true)
    try {
      isEdit ? handleEdit({
        username: info?.username || "",
        nickname: value.nickname,
        remark: value.remark
      }) : handleRegister(value)
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

  const handleRegister = async (value: MerchantRegisterReq) => {
    let { data } = await merchantRegister(value)
    params?.onOk?.();
    Modal.success({
      content: `添加成功, 密位为${data.password}`,
    });
    app.merchantList = [];
  }

  const handleEdit = async (value: MerchantUpdateReq) => {
    await merchantUpdate(value)
    params?.onOk?.();
    message.success(`修改成功`)
  }

  return (
    <>
      <Modal
        title={title}
        okText="确认"
        cancelText="取消"
        onCancel={() => {
          form.resetFields();
          params?.onCancel?.();
        }}
        confirmLoading={confirmLoading}
        open={isModalOpen}
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
          initialValues={{ ...info }}
          onFinish={onFinish}
          disabled={formDisabled}
        >
          <Form.Item<FieldType>
            name="nickname"
            label="名称"
            required
          >
            <Input />
          </Form.Item>

          <Form.Item<FieldType>
            name="remark"
            label="备注"
          >
            <TextArea rows={4} />
          </Form.Item>
        </Form >
      </Modal>
    </>
  );
};

export default MerchantCreateModal;
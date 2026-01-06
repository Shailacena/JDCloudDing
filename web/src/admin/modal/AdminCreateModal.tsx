import { useEffect, useState } from 'react';
import { Button, Form, FormProps, Input, message, Modal, Select } from 'antd';
import { AdminBaseInfoReq, IAdmin } from '../../api/types';
import { useApis } from '../../api/api';
import axios from 'axios';
import { AllRoleType, RoleType } from '../role';
import { useAppContext } from '../../AppProvider';

interface ModalDataType {
  isModalOpen: boolean
  onOk: Function;
  onCancel: Function;
  info?: FieldType
  allAdmin?: IAdmin[]
}

export type FieldType = {
  username?: string;
  role?: number;
  nickname?: string;
  remark?: string;
};

enum Title {
  CreateTxt = '新增管理员',
  EditTxt = '修改管理员'
}

const AdminCreateModal = (params: ModalDataType) => {
  const [info, setInfo] = useState(params.info)
  const [isEdit, setIsEdit] = useState(!!params.info)
  const [title, setTitle] = useState('')
  const [isModalOpen, setIsModalOpen] = useState(params.isModalOpen);
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [formDisabled, setFormDisabled] = useState<boolean>(false);
  const [roleOpts, setRoleOpts] = useState<any[]>([]);
  let { adminRegister, adminUpdate } = useApis()
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
    let o = AllRoleType.find((item) => {
      let role = RoleType.Agency
      if (app.cookie.role) {
        role = app.cookie.role
      }
      return item.roleType > role
    })
    let opts: any[] = []
    opts.push({
      value: o?.roleType,
      label: o?.label
    })
    setRoleOpts(opts)
  }, [])

  useEffect(() => {
    setTitle(isEdit ? Title.EditTxt : Title.CreateTxt)
  }, [isEdit])

  const onFinish: FormProps<AdminBaseInfoReq>['onFinish'] = async (value) => {
    setFormDisabled(true)
    setConfirmLoading(true)
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

  const handleRegister = async (value: AdminBaseInfoReq) => {
    value.createdBy = app.cookie.id!
    let { data } = await adminRegister(value)
    params?.onOk?.();
    Modal.success({
      content: `创建管理员成功, 密位为${data.password}`,
    });
  }

  const handleEdit = async (value: AdminBaseInfoReq) => {
    await adminUpdate(value)
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
            name="username"
            label="帐号"
            rules={[{ required: true, message: '请输入帐号' }]}
          >
            <Input disabled={isEdit} />
          </Form.Item>

          <Form.Item<FieldType>
            name="role"
            label="角色"
            rules={[{ required: true }]}
          >
            <Select
              disabled={isEdit}
              options={roleOpts}
            >
            </Select>
          </Form.Item>

          <Form.Item<FieldType>
            name="nickname"
            label="昵称"
            rules={[{ required: true, message: '请输入昵称' }]}
          >
            <Input />
          </Form.Item>

          {/* <Form.Item<FieldType>
            name="remark"
            label="备注"
          >
            <TextArea rows={4} />
          </Form.Item> */}
        </Form >
      </Modal>
    </>
  );
};

export default AdminCreateModal;
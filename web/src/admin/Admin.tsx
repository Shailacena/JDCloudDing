import { Space, Button, message, Card, Divider, Popconfirm, QRCode, Typography, Modal } from 'antd';
import { CaretRightOutlined, CaretDownOutlined } from '@ant-design/icons';
import { useEffect, useState } from 'react';
import { useApis } from '../api/api';
import axios from 'axios';
import AdminCreateModal, { FieldType } from './modal/AdminCreateModal';
import { getRoleName, isSuperAdmin, RoleType } from './role';
import { IAdmin } from '../api/types';
import { EnableStatus } from '../utils/constant';
import { isEnable } from '../utils/utilb';
import { convertEnable } from '../utils/biz';
import { useAppContext } from '../AppProvider';

const { Text } = Typography;

interface DataType extends IAdmin {
  key: number;
  children?: DataType[];
}

function Admin() {
  const [list, setList] = useState<DataType[]>([])
  const [isModalOpen, setIsModalOpen] = useState(false);
  let { listAdmin, adminResetPassword, adminDelete, adminEnable, adminResetVerifiCode, getMasterIncome, archiveOrders } = useApis()
  const [selectedData, setSelectedData] = useState<FieldType>(null!);
  let app = useAppContext()
  const [expandedKeys, setExpandedKeys] = useState<string[]>([]);
  // 主账号总收入状态
  const [masterTotalIncome, setMasterTotalIncome] = useState<Record<number, { totalAmount: number, loading: boolean }>>({});
  // 主账号ID集合，用于一次性加载所有主账号的收入数据
  const [masterAdminIds, setMasterAdminIds] = useState<number[]>([]);

  // 监听主账号ID集合变化，获取所有主账号的收入数据
  useEffect(() => {
    console.log('masterAdminIds', masterAdminIds)
    if (masterAdminIds.length > 0) {
      // 为每个主账号获取收入数据
      masterAdminIds.forEach(adminId => {
        // 避免重复请求
        if (!masterTotalIncome[adminId] || masterTotalIncome[adminId].loading === undefined) {
          fetchMasterIncome(adminId);
        }
      });
    }
  }, [masterAdminIds]);

  // 构建层级数据结构
  const buildHierarchicalData = (adminList: IAdmin[]) => {
    const dataMap: Record<number, DataType> = {};
    const result: DataType[] = [];
    
    // 获取当前用户信息
    const currentUserId = app?.cookie?.id;
    const currentUserRole = app?.cookie?.role;
    const isCurrentUserAdmin = currentUserRole === RoleType.Admin;
    
    // 首先创建所有节点的映射
    adminList.forEach((item, index) => {
      dataMap[item.id] = {
        ...item,
        key: index,
        children: []
      };
    });
    
    // 构建树状结构
    adminList.forEach((item) => {
      // 过滤掉当前登录用户
      if (item.id === currentUserId) {
        return;
      }
      
      const currentNode = dataMap[item.id];
      
      if (isCurrentUserAdmin && currentUserId !== undefined) {
        // 对于主管理员，只显示自己的子账号和代理账户
        if (item.parentId === currentUserId || item.masterId === currentUserId) {
          // 如果父节点是当前用户（已被过滤），则直接添加到结果中
          // 否则尝试添加到父节点
          if (item.parentId === currentUserId) {
            result.push(currentNode);
          } else {
            const parentNode = dataMap[item.parentId];
            if (parentNode && parentNode.id !== currentUserId) {
              parentNode.children?.push(currentNode);
            } else {
              // 如果父节点不存在或父节点是当前用户，直接添加到结果中
              result.push(currentNode);
            }
          }
        }
      } else {
        // 对于超级管理员或其他角色，显示完整结构
        if (item.parentId === 0 || item.parentId === currentUserId) {
          // 如果父节点是当前用户（已被过滤），则直接添加到结果中作为顶级节点
          result.push(currentNode);
        } else {
          // 子账号
          const parentNode = dataMap[item.parentId];
          if (parentNode) {
            parentNode.children?.push(currentNode);
          } else {
            // 如果父节点不存在，将当前节点添加为顶级节点
            result.push(currentNode);
          }
        }
      }
    });
    
    return result;
  };

  useEffect(() => {
    fetchListAdmin();
  }, [])

  // 获取主账号总收入
  const fetchMasterIncome = async (masterId: number) => {
    console.log('fetchMasterIncome', masterId)
    // 设置加载状态
    setMasterTotalIncome(prev => ({
      ...prev,
      [masterId]: { ...prev[masterId], loading: true }
    }));

    try {
        // 调用真实API接口获取主账号总收入
      const response = await getMasterIncome({ masterId: masterId });
      const totalAmount = response?.data?.totalIncome || 0;
      
      // 更新状态
      setMasterTotalIncome(prev => ({
        ...prev,
        [masterId]: { totalAmount, loading: false }
      }));
    } catch (error) {
      console.error('获取主账号总收入失败:', error);
      // 出错时设置默认值
      setMasterTotalIncome(prev => ({
        ...prev,
        [masterId]: { totalAmount: 0, loading: false }
      }));
    }
  };

  const fetchListAdmin = async () => {
    try {
      const { data } = await listAdmin({})
      const hierarchicalData = buildHierarchicalData(data?.list || []);
      setList(hierarchicalData);
      
      // 收集所有主账号ID（父账号ID为0的账号）
      const admins = data?.list || [];
      const masterIds = admins.filter(admin => admin.role === RoleType.Admin).map(admin => admin.id);
      setMasterAdminIds(masterIds);
    } catch (e) {
      if (axios.isAxiosError(e)) {
        let msg = e.response?.data?.message
        msg && showErrorMsg(msg);
      }
    }
  }

  const onSuccess = () => {
    fetchListAdmin()
    setIsModalOpen(false)
  };

  const openModal = (selectedData: DataType | null = null, isOpen: boolean = false) => {
    setSelectedData(selectedData!)
    setIsModalOpen(isOpen)
  }

  const showSuccessMsg = (text: string) => {
    message.success(text)
  }

  const showErrorMsg = (text: string) => {
    message.error(text)
  }

  const resetPassword = async (username: string) => {
    try {
      let { data } = await adminResetPassword({ username })
      Modal.success({
        content: `重置密码成功, 密位为 ${data.password}`,
      });
    } catch (e) {
      if (axios.isAxiosError(e)) {
        let msg = e.response?.data?.message
        msg && showErrorMsg(msg);
      }
    }
  };

  const resetVerifiCode = async (id: number) => {
    try {
      await adminResetVerifiCode({ id })
      fetchListAdmin()
      showSuccessMsg(`重置验证码成功, 查看二维码`)
    } catch (e) {
      if (axios.isAxiosError(e)) {
        let msg = e.response?.data?.message
        msg && showErrorMsg(msg);
      }
    }
  };

  const deleteAdmin = async (username: string) => {
    try {
      await adminDelete({ username })
      fetchListAdmin()
      showSuccessMsg('删除成功');
    } catch (e) {
      if (axios.isAxiosError(e)) {
        let msg = e.response?.data?.message
        msg && showErrorMsg(msg);
      }
    }
  };

  const enableAdmin = async (username: string, enable: number) => {
    try {
      await adminEnable({ username, enable })
      fetchListAdmin()
      showSuccessMsg(isEnable(enable) ? '启用成功' : '冻结成功')
    } catch (e) {
      if (axios.isAxiosError(e)) {
        let msg = e.response?.data?.message
        msg && showErrorMsg(msg);
      }
    }
  };

  // 递归渲染管理员列表
  const renderAdminList = (adminList: DataType[], level = 0) => {
    return adminList.map((admin) => {
      const isSuper = isSuperAdmin(admin.role);
      const enable = isEnable(admin.enable);
      const hasChildren = admin.children && admin.children.length > 0;
      const isExpanded = expandedKeys.includes(`admin-${admin.id}`);
      const isMaster = level === 0; // 判断是否为主账号
      
      const toggleExpand = (e: React.MouseEvent) => {
        e.stopPropagation(); // 阻止冒泡
        if (hasChildren) {
          setExpandedKeys(prev => {
            if (isExpanded) {
              return prev.filter(key => key !== `admin-${admin.id}`);
            } else {
              return [...prev, `admin-${admin.id}`];
            }
          });
        }
      };
      
      // 根据层级和角色确定背景色和边框色
      const getBackgroundColor = () => {
        if (level === 0) return '#e6f7ff'; // 主账号使用浅蓝色背景
        if (level === 1) return '#f6ffed'; // 一级子账号使用浅绿色背景
        return '#ffffff'; // 其他子账号使用白色背景
      };
      
      const getBorderColor = () => {
        if (level === 0) return '#91d5ff'; // 主账号使用蓝色边框
        if (level === 1) return '#b7eb8f'; // 一级子账号使用绿色边框
        return '#f0f0f0'; // 其他子账号使用灰色边框
      };
      
      // 获取角色标签的样式
      const getRoleTagStyle = () => {
        switch (admin.role) {
          case RoleType.SuperAdmin:
            return { backgroundColor: '#ff4d4f', color: '#fff' };
          case RoleType.Admin:
            return { backgroundColor: '#1890ff', color: '#fff' };
          case RoleType.ClonedAdmin:
            return { backgroundColor: '#52c41a', color: '#fff' };
          case RoleType.Agency:
            return { backgroundColor: '#fa8c16', color: '#fff' };
          default:
            return { backgroundColor: '#d9d9d9', color: '#333' };
        }
      };
      
      return (
        <div 
          key={admin.id} 
          style={{ 
            marginBottom: 12, 
            transition: 'all 0.3s ease',
            borderLeft: level > 0 ? `2px dashed #d9d9d9` : 'none',
            paddingLeft: level > 0 ? '18px' : '0',
            paddingTop: level > 0 ? '4px' : '0'
          }}
        >
          <div 
            style={{
              padding: '16px',
              backgroundColor: getBackgroundColor(),
              borderRadius: '6px',
              border: `1px solid ${getBorderColor()}`,
              boxShadow: isMaster ? '0 2px 8px rgba(0, 0, 0, 0.08)' : 'none',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'space-between',
              cursor: hasChildren ? 'pointer' : 'default',
              transition: 'all 0.2s ease'
            }}
            onClick={(e) => hasChildren && toggleExpand(e)}
          >
            <div style={{ display: 'flex', alignItems: 'center', flex: 1 }}>
              {hasChildren && (
                <div 
                  onClick={toggleExpand}
                  style={{ 
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    width: '24px',
                    height: '24px',
                    marginRight: '12px',
                    borderRadius: '4px',
                    backgroundColor: 'rgba(255, 255, 255, 0.8)',
                    cursor: 'pointer',
                    transition: 'all 0.2s ease'
                  }}
                  onMouseEnter={(e) => e.currentTarget.style.backgroundColor = 'rgba(255, 255, 255, 1)'}
                  onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'rgba(255, 255, 255, 0.8)'}
                >
                  {isExpanded ? <CaretDownOutlined style={{ color: '#1890ff' }} /> : <CaretRightOutlined style={{ color: '#1890ff' }} />}
                </div>
              )}
              {!hasChildren && <div style={{ width: '36px', marginRight: '12px' }} />}
              
              <div style={{ flex: 1, minWidth: 0 }}>
                <div style={{ display: 'flex', alignItems: 'center', marginBottom: '8px' }}>
                  <span 
                    style={{
                      ...getRoleTagStyle(),
                      padding: '2px 8px',
                      borderRadius: '10px',
                      fontSize: '12px',
                      fontWeight: 'normal',
                      width: '80px',
                      textAlign: 'center'
                    }}
                  >
                  {getRoleName(admin.role)}
                </span>
                  <div style={{ fontWeight: isMaster ? 'bold' : 'normal', fontSize: isMaster ? '16px' : '14px', marginLeft: '8px' }}>
                    {admin.nickname}
                  </div>
                  {hasChildren && (
                    <span style={{ marginLeft: '8px', fontSize: '12px', color: '#999' }}>
                      ({admin.children?.length}个子账号)
                    </span>
                  )}
                </div>
                <div style={{ fontSize: '12px', color: '#666', lineHeight: '1.5', display: 'flex', alignItems: 'center', flexWrap: 'wrap' }}>
                    <span style={{ marginRight: '12px' }}>账号: <strong>{admin.username}</strong></span>
                    <span style={{ marginRight: '12px' }}>状态: 
                      <span style={{ color: enable ? '#52c41a' : '#f5222d', marginLeft: '4px', fontWeight: 'bold' }}>
                        {convertEnable(admin.enable)}
                      </span>
                    </span>
                    
                    {/* 主账号总收入显示 - 与状态放在同一行 */}
                    {isMaster && (
                      <div style={{ 
                        padding: '2px 6px',
                        backgroundColor: '#f0f9ff',
                        border: '1px solid #91d5ff',
                        borderRadius: '3px',
                        display: 'inline-flex',
                        alignItems: 'center'
                      }}>
                        <span style={{ fontSize: '12px', color: '#1890ff', marginRight: '4px' }}>
                          总收入:
                        </span>
                        <span style={{ 
                          fontSize: '12px', 
                          fontWeight: 'bold', 
                          color: '#1890ff'
                        }}>
                          ¥{masterTotalIncome[admin.id]?.totalAmount?.toFixed(2) || '0.00'}
                        </span>
                        {masterTotalIncome[admin.id]?.loading && (
                          <span style={{ fontSize: '11px', color: '#1890ff', marginLeft: '4px' }}>加载中...</span>
                        )}
                      </div>
                    )}
                  </div>
              </div>
              
              <div style={{ marginLeft: '16px' }}>
                <Space size='small' wrap align="center">
                  {isMaster && admin.role === RoleType.Admin && (app?.cookie?.role === RoleType.SuperAdmin) && (
                    <Popconfirm 
                      title="归档"
                      description="归档该主账号的所有订单并删除订单"
                      onConfirm={async () => {
                        try {
                          const { data } = await archiveOrders({ adminId: admin.id })
                          showSuccessMsg(`归档成功，金额 ${data.totalAmount}，订单数 ${data.orderCount}`)
                        } catch (e) {
                          if (axios.isAxiosError(e)) {
                            let msg = e.response?.data?.message
                            msg && showErrorMsg(msg);
                          }
                        }
                      }}
                    >
                      <Button type="primary" size='small' onClick={(e) => e.stopPropagation()}>归档</Button>
                    </Popconfirm>
                  )}
                  {!isSuper && (
                    <Button
                      type="primary"
                      size='small'
                      danger={enable}
                      onClick={(e) => {
                        e.stopPropagation();
                        enableAdmin(admin.username, enable ? EnableStatus.Disabled : EnableStatus.Enabled);
                      }}
                    >
                      {enable ? '冻结' : '启用'}
                    </Button>
                  )}
                  <Button type="primary" size='small' onClick={(e) => {
                    e.stopPropagation();
                    openModal(admin, true);
                  }}>修改</Button>
                  {!isSuper && (
                    <Popconfirm 
                      title="警告" 
                      description="请确认是否删除该管理员" 
                      onConfirm={() => deleteAdmin(admin.username)}
                    >
                      <Button type="primary" size='small' danger onClick={(e) => e.stopPropagation()}>删除</Button>
                    </Popconfirm>
                  )}
                  <Button type="primary" size='small' danger onClick={(e) => {
                    e.stopPropagation();
                    resetPassword(admin.username);
                  }}>重置密码</Button>
                  <Button type="primary" size='small' onClick={(e) => {
                    e.stopPropagation();
                    resetVerifiCode(admin.id);
                  }}>重置验证码</Button>
                  <Popconfirm
                    title="验证码"
                    icon={null}
                    description={
                      admin?.urlKey ?
                        <QRCode value={admin?.urlKey} size={320} />
                        :
                        <div>
                          <Text>无二维码，点击</Text>
                          <Text type="success">重置验证码</Text>
                        </div>
                    }
                    showCancel={false}
                    okText="关闭"
                    onConfirm={() => {}}
                  >
                    <Button type="primary" size='small' onClick={(e) => e.stopPropagation()}>查看验证码</Button>
                  </Popconfirm>
                </Space>
              </div>
            </div>
          </div>
          
          {hasChildren && isExpanded && (
            <div style={{ marginTop: '8px', animation: 'fadeIn 0.3s ease-in-out' }}>
              {renderAdminList(admin.children!, level + 1)}
            </div>
          )}
        </div>
      );
    });
  };

  return (
    <>
      <Card>
        {
          (app?.cookie?.role ? app?.cookie?.role < RoleType.ClonedAdmin : false) &&
          <Button type="primary" onClick={() => { setIsModalOpen(true) }}>新增管理员</Button>
        }
        <Divider />
        <div style={{ 
          maxHeight: '80vh', 
          overflowY: 'auto',
          padding: '16px',
          backgroundColor: '#fafafa',
          borderRadius: '6px',
          border: '1px solid #e8e8e8'
        }}>
          {list.length > 0 ? (
            <div>
              {renderAdminList(list)}
            </div>
          ) : (
            <div style={{ 
              textAlign: 'center', 
              padding: '60px 20px', 
              color: '#999',
              backgroundColor: '#fff',
              borderRadius: '6px',
              border: '1px dashed #d9d9d9'
            }}>
              <div style={{ fontSize: '16px', marginBottom: '8px' }}>暂无管理员数据</div>
              <div style={{ fontSize: '12px' }}>点击上方"新增管理员"按钮添加新的管理员账号</div>
            </div>
          )}
        </div>

        {
          isModalOpen &&
          <AdminCreateModal info={selectedData} isModalOpen={isModalOpen} onOk={onSuccess} onCancel={() => openModal()} />
        }
      </Card>
    </>
  )
}

export default Admin
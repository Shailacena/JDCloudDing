export enum RoleType {
  SuperAdmin = 1,
  Admin = 2,
  ClonedAdmin = 3,
  Agency = 4
}

export interface IRoleType {
  roleType: RoleType
  label: string
}

export let AllRoleType: IRoleType[] = [
  {
    roleType: RoleType.SuperAdmin,
    label: "超级管理员",
  },
  {
    roleType: RoleType.Admin,
    label: "主账号",
  },
  {
    roleType: RoleType.ClonedAdmin,
    label: "子账号",
  },
  {
    roleType: RoleType.Agency,
    label: "代理",
  }
]

export function getRoleName(role: RoleType): string {
  for (let index = 0; index < AllRoleType.length; index++) {
      const r = AllRoleType[index];
      if (r.roleType === role) {
        return r.label
      }
    }
    return ''
}

export function isSuperAdmin(role: RoleType): boolean {
  return role === RoleType.SuperAdmin
}

export function isClonedAdmin(role: RoleType): boolean {
  return role === RoleType.ClonedAdmin
}
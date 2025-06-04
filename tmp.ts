type Department = {
  id: string;
  name: string;
  parentId: string;
  children: Department[];
};

function newDepartment(id: string, name: string, parentId: string): Department {
  return { id, name, parentId, children: [] };
}

function departmentTree(departments: Department[]): Department[] {
  const deptMap: { [id: string]: Department } = {};
  departments.forEach(dept => {
    dept.children = []; // 再利用時のクリア
    deptMap[dept.id] = dept;
  });
  departments.forEach(dept => {
    if (dept.parentId) {
      console.warn('parentId is not null', dept.parentId);
      const parent = deptMap[dept.parentId];
      if (parent) {
        parent.children.push(dept);
      }
    }
  });
  // ルートノードのみ返す
  return departments.filter(dept => !dept.parentId);
}

// 使用例
const departments = [
  newDepartment("1", "HR", ""),
  newDepartment("3", "Sales", "2"),
  newDepartment("2", "Engineering", "1"),
  newDepartment("4", "Recruitment", "6"),
  newDepartment("5", "Development", "4"),
  newDepartment("6", "Development", "5"),
];

const tree = departmentTree(departments);
console.log(JSON.stringify(tree, null, 2));
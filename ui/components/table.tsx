"use client";
import React, { useEffect, useState } from "react";
import {
    Table,
    TableHeader,
    TableColumn,
    TableBody,
    TableRow,
    TableCell,
    Input,
    Button,
    DropdownTrigger,
    Dropdown,
    DropdownMenu,
    DropdownItem,
    Chip,
    Pagination,
    Selection,
    ChipProps,
    SortDescriptor
} from "@heroui/react";
import { SearchIcon, VerticalDotsIcon, ChevronDownIcon, RefreshIcon, CopyIcon } from "@/components/icons";
import { Modal, ModalBody, ModalContent } from "@heroui/modal";
import { Card, CardBody, CardHeader } from "@heroui/card";
import axios from "axios";
import { toast } from "react-toastify";

function capitalize(str: string) {
    return str.charAt(0).toUpperCase() + str.slice(1);
}

const kindColorMap: Record<string, ChipProps["color"]> = {
    RoleBinding: "success",
    ClusterRoleBinding: "warning",
};

const columns = [
    { name: "NAME", uid: "name", sortable: true },
    { name: "KIND", uid: "kind" },
    { name: "SUBJECTS", uid: "subjects" },
    { name: "ROLE REF", uid: "role_ref" },
    { name: "DETAILS", uid: "details" },
];

const kindOptions = [
    { name: "ClusterRoleBinding", uid: "ClusterRoleBinding" },
    { name: "RoleBinding", uid: "RoleBinding" },
    { name: "ProjectRoleTemplateBinding", uid: "projectroletemplatebindings" },
    { name: "ClusterRoleTemplateBinding", uid: "clusterroletemplatebindings" },
    { name: "GlobalRoleBinding", uid: "globalrolebindings" },
];

type Subject = {
    kind: string;
    apiGroup: string;
    name: string;
};

type RoleRef = {
    kind: string;
    apiGroup: string;
    name: string;
};

type BindingData = {
    id: number;
    name: string;
    kind: string;
    subjects: Subject[];
    roleRef: RoleRef;
    raw?: string;
};

export default function MainTable() {
    const [data, setData] = useState<BindingData[]>([]);
    const [filterValue, setFilterValue] = React.useState("");
    const [selectedKeys, setSelectedKeys] = React.useState<Selection>(new Set([]));
    const [visibleColumns, setVisibleColumns] = React.useState<Selection>(new Set(columns.map(column => column.uid)));
    const [kindFilter, setKindFilter] = React.useState<Selection>("all");
    const [rowsPerPage, setRowsPerPage] = React.useState(5);
    const [sortDescriptor, setSortDescriptor] = React.useState<SortDescriptor>({
        column: "name",
        direction: "ascending",
    });
    const [isModalOpen, setIsModalOpen] = React.useState(false);
    const [modalData, setModalData] = React.useState<BindingData | any | null>(null);
    const [page, setPage] = React.useState(1);
    const hasSearchFilter = Boolean(filterValue);

    useEffect(() => {
        axios.get('/api/data')
            .then(response => setData(response.data))
            .catch(error => console.error('Error fetching data:', error));
    }, []);

    useEffect(() => {
        const handleKeyDown = (event: KeyboardEvent) => {
            if (event.key === "Escape") {
                setIsModalOpen(false);
            }
        };

        document.addEventListener("keydown", handleKeyDown);

        return () => {
            document.removeEventListener("keydown", handleKeyDown);
        };
    }, []);

    const headerColumns = React.useMemo(() => {
        if (visibleColumns === "all") return columns;

        return columns.filter((column) => Array.from(visibleColumns).includes(column.uid));
    }, [visibleColumns]);

    const filteredItems = React.useMemo(() => {
        let filteredUsers = [...data];

        if (hasSearchFilter) {
            filteredUsers = filteredUsers.filter((user) =>
                user.name.toLowerCase().includes(filterValue.toLowerCase()),
            );
        }
        if (kindFilter !== "all" && Array.from(kindFilter).length !== kindOptions.length) {
            filteredUsers = filteredUsers.filter((user) =>
                Array.from(kindFilter).includes(user.kind),
            );
        }

        return filteredUsers;
    }, [data, filterValue, hasSearchFilter, kindFilter]);

    const pages = Math.ceil(filteredItems.length / rowsPerPage);

    const items = React.useMemo(() => {
        const start = (page - 1) * rowsPerPage;
        const end = start + rowsPerPage;

        return filteredItems.slice(start, end);
    }, [page, filteredItems, rowsPerPage]);

    const sortedItems = React.useMemo(() => {
        return [...items].sort((a: BindingData, b: BindingData) => {
            const first = a[sortDescriptor.column as keyof BindingData] as string | number;
            const second = b[sortDescriptor.column as keyof BindingData] as string | number;
            const cmp = first < second ? -1 : first > second ? 1 : 0;

            return sortDescriptor.direction === "descending" ? -cmp : cmp;
        });
    }, [sortDescriptor, items]);

    const renderCell = React.useCallback((data: BindingData, columnKey: React.Key) => {
        const cellValue = data[columnKey as keyof BindingData];

        switch (columnKey) {
            case "name":
                return (
                    <div className="flex flex-col">
                        <p className="text-bold text-small">{cellValue?.toString()}</p>
                        <p className="text-bold text-tiny text-default-400">{data.name}</p>
                    </div>
                );
            case "kind":
                return (
                    <Chip color={kindColorMap[data.kind]} size="sm" variant="flat">
                        {cellValue?.toString()}
                    </Chip>
                );
            case "subjects":
                return (
                    <div>
                        {data.subjects?.map((subject: Subject, index) => (
                            <p key={index}>{subject.kind} - {subject.name}</p>
                        ))}
                    </div>
                );
            case "role_ref":
                return (
                    <div>
                        <p>{data.roleRef?.kind} - {data.roleRef?.apiGroup} - {data.roleRef?.name}</p>
                    </div>
                );
            case "details":
                return (
                    <div className="relative flex justify-center items-center gap-2">
                        <Dropdown>
                            <DropdownTrigger aria-label="More options">
                                <Button isIconOnly size="sm" variant="light">
                                    <VerticalDotsIcon className="text-default-300" />
                                </Button>
                            </DropdownTrigger>
                            <DropdownMenu aria-label="Details options">
                                <DropdownItem 
                                    key="view"
                                    onClick={
                                        () => {
                                            setModalData(data); // Set the binding data to the modalData state
                                            setIsModalOpen(true); // Open the modal
                                        }
                                    }>View</DropdownItem>
                            </DropdownMenu>
                        </Dropdown>
                    </div>
                );
            default:
                return typeof cellValue === 'string' || typeof cellValue === 'number' ? cellValue : JSON.stringify(cellValue);
        }
    }, []);

    const onNextPage = React.useCallback(() => {
        if (page < pages) {
            setPage(page + 1);
        }
    }, [page, pages]);

    const onPreviousPage = React.useCallback(() => {
        if (page > 1) {
            setPage(page - 1);
        }
    }, [page]);

    const onRowsPerPageChange = React.useCallback((e: React.ChangeEvent<HTMLSelectElement>) => {
        setRowsPerPage(Number(e.target.value));
        setPage(1);
    }, []);

    const onSearchChange = React.useCallback((value?: string) => {
        if (value) {
            setFilterValue(value);
            setPage(1);
        } else {
            setFilterValue("");
        }
    }, []);

    const onClear = React.useCallback(() => {
        setFilterValue("")
        setPage(1)
    }, [])

    const topContent = React.useMemo(() => {
        return (
            <div className="flex flex-col gap-4">
                <div className="flex justify-between gap-3 items-end">
                    <Input
                        isClearable
                        className="w-full sm:max-w-[44%]"
                        placeholder="Search by name..."
                        startContent={<SearchIcon />}
                        value={filterValue}
                        onClear={() => onClear()}
                        onValueChange={onSearchChange}
                    />
                    <div className="flex gap-3">
                        <Dropdown>
                            <DropdownTrigger className="hidden sm:flex" aria-label="Filter by kind">
                                <Button endContent={<ChevronDownIcon className="text-small" />} variant="flat">
                                    Kind
                                </Button>
                            </DropdownTrigger>
                            <DropdownMenu
                                disallowEmptySelection
                                aria-label="Kind filter options"
                                closeOnSelect={false}
                                selectedKeys={kindFilter}
                                selectionMode="multiple"
                                onSelectionChange={setKindFilter}
                            >
                                {kindOptions.map((kind) => (
                                    <DropdownItem key={kind.uid} className="capitalize">
                                        {capitalize(kind.name)}
                                    </DropdownItem>
                                ))}
                            </DropdownMenu>
                        </Dropdown>
                        <Dropdown>
                            <DropdownTrigger className="hidden sm:flex" aria-label="Select columns">
                                <Button endContent={<ChevronDownIcon className="text-small" />} variant="flat">
                                    Columns
                                </Button>
                            </DropdownTrigger>
                            <DropdownMenu
                                disallowEmptySelection
                                aria-label="Column selection options"
                                closeOnSelect={false}
                                selectedKeys={visibleColumns}
                                selectionMode="multiple"
                                onSelectionChange={setVisibleColumns}
                            >
                                {columns.map((column) => (
                                    <DropdownItem key={column.uid} className="capitalize">
                                        {capitalize(column.name)}
                                    </DropdownItem>
                                ))}
                            </DropdownMenu>
                        </Dropdown>
                        <Button color="primary" endContent={<RefreshIcon />} onClick={() => axios.get('/api/data').then(response => setData(response.data)).catch(error => console.error('Error fetching data:', error))}>
                            Refresh
                        </Button>
                    </div>
                </div>
                <div className="flex justify-between items-center">
                    <span className="text-default-400 text-small">Total {data.length} bindings</span>
                    <label className="flex items-center text-default-400 text-small">
                        Rows per page:
                        <select
                            className="bg-transparent outline-none text-default-400 text-small"
                            onChange={onRowsPerPageChange}
                        >
                            <option value="5">5</option>
                            <option value="10">10</option>
                            <option value="15">15</option>
                        </select>
                    </label>
                </div>
            </div>
        );
    }, [filterValue, onSearchChange, kindFilter, visibleColumns, onRowsPerPageChange, onClear, data]);

    const bottomContent = React.useMemo(() => {
        return (
            <div className="py-2 px-2 flex justify-between items-center">
                <Pagination
                    isCompact
                    showControls
                    showShadow
                    color="primary"
                    page={page}
                    total={pages}
                    onChange={setPage}
                />
                <div className="hidden sm:flex w-[30%] justify-end gap-2">
                    <Button isDisabled={pages === 1} size="sm" variant="flat" onPress={onPreviousPage}>
                        Previous
                    </Button>
                    <Button isDisabled={pages === 1} size="sm" variant="flat" onPress={onNextPage}>
                        Next
                    </Button>
                </div>
            </div>
        );
    }, [page, pages, onPreviousPage, onNextPage]);

    const copyToClipboard = async () => {
        if (modalData && modalData.raw) {
            try {
                await navigator.clipboard.writeText(modalData.raw);
                toast.success("Successfully copied to the clipboard!");
            } catch (err) {
                toast.error("Failed to copy to the clipboard!");
                console.error("Failed to copy to the clipboard:", err)
            }
        }
    };

    return (
        <>
            <Card
                className="p-2"
            >
                <CardHeader>
                    <h2>RBAC Table</h2>
                </CardHeader>
                <CardBody>
                    <Table
                        aria-label="Example table with custom cells, pagination and sorting"
                        isHeaderSticky
                        isStriped
                        bottomContent={bottomContent}
                        bottomContentPlacement="outside"
                        classNames={{
                            wrapper: "max-h-[382px]",
                        }}
                        selectedKeys={selectedKeys}
                        selectionMode="none"
                        sortDescriptor={sortDescriptor}
                        topContent={topContent}
                        topContentPlacement="outside"
                        onSelectionChange={setSelectedKeys}
                        onSortChange={setSortDescriptor}
                    >
                        <TableHeader columns={headerColumns}>
                            {(column) => (
                                <TableColumn
                                    key={column.uid}
                                    align={column.uid === "details" ? "center" : "start"}
                                    allowsSorting={column.sortable}
                                >
                                    {column.name}
                                </TableColumn>
                            )}
                        </TableHeader>
                        <TableBody emptyContent={"No bindings found"} items={sortedItems}>
                            {(item) => (
                                <TableRow>
                                    {(columnKey) => <TableCell>{renderCell(item, columnKey)}</TableCell>}
                                </TableRow>
                            )}
                        </TableBody>
                    </Table>
                    {/* Modal to display the data */}
                    <Modal
                        size="3xl"
                        radius="md"
                        shadow="lg"
                        motionProps={{
                            variants: {
                                enter: {
                                    y: 0,
                                    opacity: 1,
                                    transition: {
                                        duration: 0.3,
                                        ease: "easeOut",
                                    },
                                },
                                exit: {
                                    y: -20,
                                    opacity: 0,
                                    transition: {
                                        duration: 0.2,
                                        ease: "easeIn",
                                    },
                                },
                            }}}
                        isOpen={isModalOpen}
                        onClose={() => setIsModalOpen(false)}
                    >
                        <ModalContent>
                            <ModalBody>
                                <Card className="p-2 m-3" isBlurred shadow="sm" style={{ maxHeight: '70vh', overflow: 'auto' }}>
                                    <CardBody className="relative">
                                        <div className="absolute top-0 right-0 mb-2">
                                            <Button isIconOnly size="sm" variant="light" aria-label="Copy data" onClick={copyToClipboard}>
                                                <CopyIcon className="text-default-300" />
                                            </Button>
                                        </div>
                                            <pre style={{ whiteSpace: 'pre-wrap' }}>
                                                {modalData && modalData.raw}
                                            </pre>
                                    </CardBody>
                                </Card>
                            </ModalBody>
                        </ModalContent>
                    </Modal>
                </CardBody>
            </Card>
        </>
    );
}
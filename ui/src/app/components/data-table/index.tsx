import { TableCell } from '@/app/components/base/tables/table-cell';
import { TableRow } from '@/app/components/base/tables/table-row';
import {
  useState,
  useCallback,
  useRef,
  useEffect,
  FC,
  HTMLAttributes,
} from 'react';

interface ScrollableResizableTableProps
  extends HTMLAttributes<HTMLTableElement> {
  clms: { name: string; key: string; width?: number }[];
  isActionable?: boolean;
  isExpandable?: boolean;
  ontoggle?: (boolean) => void;
  isOptionable?: boolean;
}

export const ScrollableResizableTable: FC<ScrollableResizableTableProps> = ({
  clms,
  isActionable = true,
  isExpandable = false,
  ontoggle,
  children,
  isOptionable,
}) => {
  const [columns, setColumns] = useState<
    { key: string; width: number; name: string }[]
  >([]);

  useEffect(() => {
    setColumns(
      clms.map(x => ({
        name: x.name,
        key: x.key,
        width: 200,
      })),
    );
  }, [clms]);
  const tableRef = useRef(null);
  const [tableWidth, setTableWidth] = useState(0);

  useEffect(() => {
    const updateTableWidth = () => {
      if (tableRef.current) {
        const newWidth = columns.reduce((sum, column) => sum + column.width, 0);
        setTableWidth(newWidth);
      }
    };
    updateTableWidth();
  }, [columns]);

  const handleResize = useCallback((index, newWidth) => {
    setColumns(prevColumns =>
      prevColumns.map((column, i) =>
        i === index ? { ...column, width: Math.max(100, newWidth) } : column,
      ),
    );
  }, []);

  return (
    <div className="w-full flex flex-col flex-1 overflow-x-auto">
      <div
        ref={tableRef}
        style={{
          minWidth: `${tableWidth}px`,
          //   minWidth: '100%',
          maxWidth: '100% !important',
        }}
      >
        <table className="w-full border-collapse bg-white dark:bg-gray-900 border-y ">
          <thead className="uppercase">
            <TableRow className="bg-gray-100 dark:bg-gray-950">
              {isActionable && (
                <TableCell className="w-8 h-8 ">
                  <div className="w-8 h-8 flex justify-center items-center">
                    <input
                      type="checkbox"
                      onChange={x => ontoggle && ontoggle(x.target.checked)}
                    />
                  </div>
                </TableCell>
              )}

              {isExpandable && (
                <TableCell className="w-8 h-8 ">
                  <div className="w-8 h-8 flex justify-center items-center"></div>
                </TableCell>
              )}

              {columns.map((column, index) => (
                <TableCell
                  key={column.key}
                  className="px-2 py-2 text-left uppercase text-xs/6 font-medium tracking-wider relative text-muted"
                  style={{ width: column.width }}
                >
                  {column.name}
                  {index !== columns.length - 1 && (
                    <div
                      className="absolute top-1 -right-[2px] bottom-1 w-[2px] cursor-col-resize hover:bg-blue-500"
                      onMouseDown={e => {
                        e.preventDefault();
                        const startX = e.pageX;
                        const startWidth = column.width;

                        const onMouseMove = e => {
                          const newWidth = startWidth + e.pageX - startX;
                          handleResize(index, newWidth);
                        };

                        const onMouseUp = () => {
                          document.removeEventListener(
                            'mousemove',
                            onMouseMove,
                          );
                          document.removeEventListener('mouseup', onMouseUp);
                        };

                        document.addEventListener('mousemove', onMouseMove);
                        document.addEventListener('mouseup', onMouseUp);
                      }}
                    ></div>
                  )}
                </TableCell>
              ))}

              {isOptionable && <TableCell className="w-10 h-10"></TableCell>}
            </TableRow>
          </thead>
          <tbody>{children}</tbody>
        </table>
      </div>
    </div>
  );
};

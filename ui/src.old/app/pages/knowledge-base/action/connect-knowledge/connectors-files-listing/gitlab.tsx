import { FC, memo, useCallback, useContext, useEffect, useState } from 'react';
import { cn } from '@/utils';
import { BluredWrapper } from '@/app/components/wrapper/blured-wrapper';
import { SearchIconInput } from '@/app/components/form/input/IconInput';
import { Card } from '@/app/components/base/cards';
import { ConnectorFileContext } from '@/hooks/use-connector-file-page-store';
import { TablePagination } from '@/app/components/base/tables/table-pagination';
import { useCredential, useRapidaStore } from '@/hooks';
import { Struct } from 'google-protobuf/google/protobuf/struct_pb';
import toast from 'react-hot-toast/headless';
import { Spinner } from '@/app/components/loader/spinner';
import { KnowledgeFileListingProps } from '@/app/pages/knowledge-base/action/connect-knowledge/connectors-files-listing';
import { TableHederWithCheckbox } from '@/app/components/base/tables/table-head';
import { TableRow } from '@/app/components/base/tables/table-row';
import { TableCell } from '@/app/components/base/tables/table-cell';
import { SizeCell } from '@/app/components/base/tables/size-cell';
import { TextCell } from '@/app/components/base/tables/text-cell';

export const GitlabKnowledgeFileListing: FC<KnowledgeFileListingProps> = memo(
  ({ toolProvider, className }) => {
    const ctx = useContext(ConnectorFileContext);
    const [userId, token, projectId] = useCredential();
    const { loading, showLoader, hideLoader } = useRapidaStore();
    const [selectedFiles, setSelectedFiles] = useState<string[]>([]);

    const [active, setActive] = useState('all');

    const ontoggle = id => {
      let allWithoutCurrent = selectedFiles.filter(x => x !== id);
      if (allWithoutCurrent.length === selectedFiles.length) {
        setSelectedFiles([...selectedFiles, id]);
      } else {
        setSelectedFiles(allWithoutCurrent);
      }
    };

    const ontoggleall = check => {
      if (check) {
        setSelectedFiles(
          ctx.filterFiles.map(
            x => x.getFieldsMap().get('id')?.getStringValue()!,
          ),
        );
        return;
      }
      setSelectedFiles([]);
    };

    /**
     *
     */
    const onSuccess = useCallback((s: Struct[]) => {
      hideLoader();
      ctx.onChangeFiles(s);
    }, []);

    /**
     *
     */
    const onError = useCallback((err: string) => {
      hideLoader();
      toast.error(err);
    }, []);

    /**
     *
     */
    useEffect(() => {
      showLoader();
      ctx.getAllConnectorFiles(
        toolProvider.getId(),
        token,
        userId,
        projectId,
        onError,
        onSuccess,
      );
    }, [toolProvider, ctx.page, ctx.pageSize, JSON.stringify(ctx.criteria)]);

    return (
      <Card className={cn('overflow-auto relative p-0', className)}>
        <BluredWrapper
          className={cn(
            'border-none sticky top-0 z-1 dark:bg-gray-950 flex-col items-start',
          )}
        >
          <div className="flex w-full">
            <SearchIconInput
              iconClassName="w-4 h-4"
              className="pl-7"
              onChange={t => {
                ctx.addCriteria('title', t.target.value, '=');
              }}
            />
            <TablePagination
              currentPage={ctx.page}
              onChangeCurrentPage={ctx.setPage}
              totalItem={ctx.totalCount}
              pageSize={ctx.pageSize}
              onChangePageSize={ctx.setPageSize}
            />
          </div>
          <ul className="flex flex-wrap pt-4">
            <li className="m-1">
              <button
                className={cn(
                  'inline-flex items-center justify-center text-sm rounded-[2px] border border-transparent duration-150 font-medium ease-in-out leading-5 shadow-sm px-3 py-1',
                  'capitalize',
                  'all' === active && 'bg-blue-500! text-white!',
                  'border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-gray-500 dark:text-gray-400 dark:hover:border-gray-600 hover:border-gray-300',
                )}
                onClick={() => {
                  setActive('all');
                  ctx.addCriteria('mimeType', '', '!=');
                }}
              >
                All
              </button>
            </li>
            <li className="m-1">
              <button
                className={cn(
                  'inline-flex items-center justify-center text-sm rounded-[2px] border border-transparent duration-150 font-medium ease-in-out leading-5 shadow-sm px-3 py-1',
                  'capitalize',
                  'directories' === active && 'bg-blue-500! text-white!',
                  'border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-gray-500 dark:text-gray-400 dark:hover:border-gray-600 hover:border-gray-300',
                )}
                onClick={() => {
                  setActive('directories');
                  ctx.addCriteria(
                    'mimeType',
                    'application/vnd.google-apps.folder',
                    '=',
                  );
                }}
              >
                Directories
              </button>
            </li>
            <li className="m-1">
              <button
                className={cn(
                  'inline-flex items-center justify-center text-sm rounded-[2px] border border-transparent duration-150 font-medium ease-in-out leading-5 shadow-sm px-3 py-1',
                  'capitalize',
                  'files' === active && 'bg-blue-500! text-white!',
                  'border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-gray-500 dark:text-gray-400 dark:hover:border-gray-600 hover:border-gray-300',
                )}
                onClick={() => {
                  setActive('files');
                  ctx.addCriteria(
                    'mimeType',
                    'application/vnd.google-apps.folder',
                    '!=',
                  );
                }}
              >
                Files
              </button>
            </li>
          </ul>
        </BluredWrapper>

        <table className="text-sm w-full table-fixed">
          <TableHederWithCheckbox
            ontoggle={ontoggleall}
            columns={[
              {
                name: 'Name',
                key: 'fileName',
              },
              {
                name: 'Folder',
                key: 'folderName',
              },
              {
                name: 'size',
                key: 'fileSize',
              },
              {
                name: 'type',
                key: 'fileType',
              },
            ]}
          />
          <tbody>
            {ctx.files.map((x, idx) => {
              return (
                <TableRow key={idx}>
                  <TableCell>
                    <input
                      type="checkbox"
                      name="file-ids"
                      value={x.getFieldsMap().get('id')?.getStringValue()}
                      checked={selectedFiles.some(
                        y => y === x.getFieldsMap().get('id')?.getStringValue(),
                      )}
                      onChange={e =>
                        ontoggle(x.getFieldsMap().get('id')?.getStringValue())
                      }
                    />
                  </TableCell>
                  <TextCell>
                    {x.getFieldsMap().get('title')?.getStringValue()}
                  </TextCell>
                  <TableCell>
                    {x.getFieldsMap().get('folder')?.getStringValue()}
                  </TableCell>
                  <SizeCell
                    size={x.getFieldsMap().get('fileSize')?.getStringValue()}
                  />
                  <TableCell>
                    {x.getFieldsMap().get('mimeType')?.getStringValue()}
                  </TableCell>
                </TableRow>
              );
            })}
          </tbody>
        </table>

        {loading ? (
          <div className="py-8 flex justify-center flex-col items-center">
            <Spinner size="md" />
          </div>
        ) : (
          <>
            {ctx.files.length === 0 && (
              <div className="py-8 flex justify-center flex-col items-center">
                <h2 className="font-bold text-lg">No files or folders</h2>
                <h3 className="text-base font-medium">
                  We are not able to find any files are folder in given
                  dataseouce.
                </h3>
              </div>
            )}
          </>
        )}
      </Card>
    );
  },
);

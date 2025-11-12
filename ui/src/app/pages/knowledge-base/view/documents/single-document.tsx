import { FC } from 'react';
import { TickIcon } from '@/app/components/Icon/Tick';
import { cn } from '@/utils';
import { toHumanReadableRelativeTime } from '@/utils/date';
import { KnowledgeDocument } from '@rapidaai/react';
import { useKnowledgeDocumentPageStore } from '@/hooks/use-knowledge-document-page-store';
import { formatFileSize, formatNumber } from '@/utils/format';
import { DocumentSourcePill } from '@/app/components/pill/document-source-pill';
import { FileExtensionIcon } from '@/app/components/Icon/file-extension';
import { DocumentOption } from '@/app/pages/knowledge-base/view/documents/document-option';
import { useCredential, useRapidaStore } from '@/hooks';
import toast from 'react-hot-toast/headless';
import { WarningInfo } from '@/app/components/Icon/Warning';
import { TableRow } from '@/app/components/base/tables/table-row';
import { TableCell } from '@/app/components/base/tables/table-cell';
import { LabelCell } from '@/app/components/base/tables/label-cell';

/**
 *
 */
interface SingleDocumentProps {
  /**
   * current endpoint
   */
  document: KnowledgeDocument;

  /**
   *
   */
  onReload: () => void;
}

/**
 *
 * @param props
 * @returns
 */
export const SingleDocument: FC<SingleDocumentProps> = ({
  document,
  onReload,
}) => {
  const kdAction = useKnowledgeDocumentPageStore();
  const [userId, token, projectId] = useCredential();
  const { showLoader, hideLoader } = useRapidaStore();
  const onerror = (err: string) => {
    hideLoader();
    toast.error(err);
    onReload();
  };
  const onsuccess = (e: boolean) => {
    hideLoader();
    onReload();
  };

  const onReloadIndex = (
    knowledgeId: string,
    knowledgeDocumentId: string[],
    indexType: string,
  ) => {
    showLoader();
    kdAction.indexKnowledgeDocument(
      knowledgeId,
      knowledgeDocumentId,
      indexType,
      projectId,
      token,
      userId,
      onerror,
      onsuccess,
    );
  };
  return (
    <TableRow
      data-id={`doc-${document.getId()}`}
      x-knowledge-id={`knowledge-${document.getKnowledgeid()}`}
    >
      {kdAction.visibleColumn('getStatus') && (
        <TableCell>
          <div className="flex items-center space-x-1.5">
            <span
              className={cn(
                document.getDisplaystatus() === 'available' &&
                  'text-green-600! bg-green-400/20!',
                document.getDisplaystatus() === 'error' &&
                  'text-rose-600! bg-rose-400/20!',
                'p-1 bg-yellow-400/20 text-yellow-600 rounded-[2px] w-fit block',
              )}
            >
              {document.getDisplaystatus() === 'error' ? (
                <WarningInfo className="w-6 h-6" />
              ) : (
                <TickIcon className="w-6 h-6" />
              )}
            </span>
            <div>
              <span
                className={cn(
                  document.getDisplaystatus() === 'available' &&
                    '!text-green-60',
                  document.getDisplaystatus() === 'error' && 'text-rose-600!',
                  'font-medium block leading-3 capitalize',
                )}
              >
                {document.getDisplaystatus()}
              </span>
              <span className="opacity-60 text-xs leading-3 capitalize truncate">
                {document.getDisplaystatus()}{' '}
                {document?.getCreateddate() &&
                  toHumanReadableRelativeTime(document?.getCreateddate()!)}
              </span>
            </div>
          </div>
        </TableCell>
      )}
      {kdAction.visibleColumn('getName') && (
        <td className="px-2 py-2 text-left text-sm font-medium my-auto relative w-auto flex space-x-1 items-center">
          <div className="p-1.5 border rounded-[2px] bg-gray-50 dark:bg-gray-950/30 backdrop-blur-sm mr-1">
            <FileExtensionIcon filename={document.getName()} />
          </div>
          <div className="flex flex-col grow flex-1">
            <span className="font-semibold line-clamp-1 text-[0.9rem]">
              {document.getName()}
            </span>
            <span className="font-medium truncate text-[0.8rem] opacity-75">
              Uploaded on{' '}
              {document.getCreateddate() &&
                toHumanReadableRelativeTime(document.getCreateddate()!)}
            </span>
          </div>
        </td>
      )}
      {kdAction.visibleColumn('getDocumenttype') && (
        <LabelCell className="bg-blue-300/10 text-blue-500 dark:text-blue-400 truncate">
          {document
            .getDocumentsource()
            ?.getFieldsMap()
            .get('mimeType')
            ?.getStringValue()}
        </LabelCell>
      )}

      {kdAction.visibleColumn('getDocumentSource') && (
        <TableCell>
          <DocumentSourcePill
            type={document
              .getDocumentsource()
              ?.getFieldsMap()
              .get('type')
              ?.getStringValue()}
            source={document
              .getDocumentsource()
              ?.getFieldsMap()
              .get('source')
              ?.getStringValue()}
          />
        </TableCell>
      )}

      {kdAction.visibleColumn('getDocumentsize') && (
        <LabelCell>{formatFileSize(document.getDocumentsize())}</LabelCell>
      )}

      {kdAction.visibleColumn('getRetrievalcount') && (
        <LabelCell>{formatFileSize(document.getRetrievalcount())}</LabelCell>
      )}

      {kdAction.visibleColumn('getTokencount') && (
        <LabelCell>{formatNumber(document.getTokencount())}</LabelCell>
      )}

      {kdAction.visibleColumn('getWordcount') && (
        <LabelCell>{formatNumber(document.getWordcount())}</LabelCell>
      )}
      {kdAction.visibleColumn('getId') && (
        <LabelCell>{`doc_${document.getId()}`}</LabelCell>
      )}
      <TableCell>
        <DocumentOption document={document} onReloadIndex={onReloadIndex} />
      </TableCell>
    </TableRow>
  );
};

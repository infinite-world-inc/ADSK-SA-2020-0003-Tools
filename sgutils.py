"""
DreamView 2021
Author: Mark Thielen
Date: May 12, 2021

"""
import datetime
import json
import os

source_path = r'D:\shotgun\source_files'
fix_path = r'D:\shotgun\fix_phage'

def process(sg):
    """process shotgun files

    Args:
        sg (:obj:): shotgun instance
    """
    print('testing from process...')
    tags = [{'type': 'Tag', 'id': 4379}]
    any_filters = [
        ['filename', 'ends_with', '.mb'],
        ['filename', 'ends_with', '.ma']]

    filters = [
        # ['id', 'in', [1717073, 1659356, 1073338]],
        ['created_at', 'greater_than', datetime.datetime(2021, 4, 1)],
        ['tags', 'not_in', tags],
        {'filter_operator': 'any',
                            'filters': any_filters}
    ]
    fields = list(sg.schema_field_read('Attachment').keys())
    # fields = ['filename', 'created_at', 'attachment_links', 'project']
    attachments = sg.find('Attachment', filters, fields)

    source_files = {}
    bad = []
    not_bad = []
    for attachment in attachments:
        print('Processing: {}'.format(attachment.get('filename')))
        path = _download_file(sg, attachment)

        if not path:
            continue

        suspect = open(path, 'rb').read()
        if b'phage' in suspect:
            print('Bad: {}'.format(path))
            bad.append(path)
            attachment['source_file'] = path
            attachment['source_file_id'] = attachment['id']
            attachment['created_at'] = None
            attachment['updated_at'] = None
            source_files[os.path.basename(path)] = attachment

            with open(os.path.join(source_path, 'source_files.json'), 'a') as f:
                f.write(json.dumps(source_files, indent=4))
            sg.update("Attachment", attachment['id'], {'tags': tags})
        else:
            not_bad.append(path)
            try:
                os.remove(path)
            except IOError as e:
                print(e)

def _download_file(sg, attachment):
    date = attachment['created_at']

    source_file = 'source_{}{:02d}{:02d}_{}_{}'.format(date.year, date.month, date.day, attachment['id'], attachment.get('filename'))
    source_path_file = os.path.join(source_path, source_file)

    # attachment = {'type': 'Attachment', 'id': attachment_id }
    if os.path.isfile(source_path_file):
        return

    sg.download_attachment(attachment, file_path=source_path_file)

    return source_path_file


if __name__ == '__main__':
    process(sg)
